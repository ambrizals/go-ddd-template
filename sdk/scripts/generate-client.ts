import { readFileSync, writeFileSync, mkdirSync, existsSync } from "fs";
import { resolve, dirname } from "path";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));

interface SwaggerSpec {
  paths: Record<string, Record<string, OperationObject>>;
  definitions?: Record<string, SchemaObject>;
  basePath?: string;
}

interface OperationObject {
  tags?: string[];
  summary?: string;
  parameters?: ParameterObject[];
  responses?: Record<string, ResponseObject>;
}

interface ParameterObject {
  name: string;
  in: string;
  required?: boolean;
  schema?: SchemaObject;
}

interface ResponseObject {
  schema?: SchemaObject;
  content?: Record<string, SchemaObject>;
}

interface SchemaObject {
  type?: string;
  properties?: Record<string, SchemaObject>;
  $ref?: string;
  items?: SchemaObject;
}

function getRefName(ref: string): string {
  const parts = ref.split("/");
  return parts[parts.length - 1];
}

function extractType(
  schema: SchemaObject,
  definitions: Record<string, SchemaObject>,
): string {
  if (!schema) return "unknown";

  if (schema.$ref) {
    return getRefName(schema.$ref);
  }

  if (schema.type === "integer") return "number";
  if (schema.type === "array" && schema.items) {
    return `${extractType(schema.items, definitions)}[]`;
  }
  if (schema.type === "object" || schema.properties) {
    const props = schema.properties
      ? Object.entries(schema.properties).map(
          ([key, val]) => `${key}: ${extractType(val, definitions)}`,
        )
      : [];
    return props.length > 0
      ? `{ ${props.join("; ")} }`
      : "Record<string, unknown>";
  }

  return schema.type || "unknown";
}

function getResourceName(path: string): string {
  const segments = path.split("/").filter(Boolean);
  const lastSegment = segments[segments.length - 1];

  if (lastSegment.includes("{id}")) {
    return segments[segments.length - 2] || "resource";
  }
  return lastSegment || "resource";
}

function getMethodName(path: string, method: string): string {
  const resource = getResourceName(path);

  switch (method) {
    case "get":
      if (path.includes("{id}")) return "getById";
      return "getAll";
    case "post":
      return "create";
    case "put":
    case "patch":
      return "update";
    case "delete":
      return "delete";
    default:
      return method;
  }
}

function getBodyType(
  op: OperationObject,
  definitions: Record<string, SchemaObject>,
): string | null {
  const bodyParam = op.parameters?.find((p) => p.in === "body");
  if (bodyParam?.schema) {
    if (bodyParam.schema.$ref) {
      return `Schemas["${getRefName(bodyParam.schema.$ref)}"]`;
    }
    return extractType(bodyParam.schema, definitions);
  }
  return null;
}

function getPathParams(op: OperationObject): { name: string; type: string }[] {
  return (op.parameters?.filter((p) => p.in === "path") || []).map((p) => ({
    name: p.name,
    type: p.schema?.type === "integer" ? "number" : "string",
  }));
}

function getResponseType(
  op: OperationObject,
  definitions: Record<string, SchemaObject>,
): string {
  const successResponse = Object.entries(op.responses || {}).find(([code]) =>
    code.startsWith("2"),
  );
  if (!successResponse) return "unknown";

  const [_, response] = successResponse;

  if (response.schema?.$ref) {
    return `Schemas["${getRefName(response.schema.$ref)}"]`;
  }

  if (response.content) {
    const jsonContent =
      response.content["application/json"] || response.content["*/*"];
    if (jsonContent?.$ref) {
      return `Schemas["${getRefName(jsonContent.$ref)}"]`;
    }
    if (jsonContent) {
      return extractType(jsonContent, definitions);
    }
  }
  return "unknown";
}

function generateClient(spec: SwaggerSpec): string {
  const definitions = spec.definitions || {};
  const basePath = spec.basePath || "/api/v1";
  const paths = spec.paths || {};

  const resources: Record<
    string,
    { methods: { path: string; httpMethod: string; op: OperationObject }[] }
  > = {};

  for (const [routePath, methods] of Object.entries(paths)) {
    const resourceName = getResourceName(routePath);

    if (!resources[resourceName]) {
      resources[resourceName] = { methods: [] };
    }

    for (const [httpMethod, op] of Object.entries(methods)) {
      if (["get", "post", "put", "patch", "delete"].includes(httpMethod)) {
        resources[resourceName].methods.push({
          path: routePath,
          httpMethod,
          op,
        });
      }
    }
  }

  let code = `// Auto-generated SDK Client
import type { paths, components } from "./api";

type Schemas = components["schemas"];

export type ClientOptions = {
  baseUrl: string;
  headers?: Record<string, string>;
};

export class SDKError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: unknown
  ) {
    super(message);
    this.name = "SDKError";
  }
}

async function request<T>(
  url: string,
  options: RequestInit,
  baseUrl: string
): Promise<T> {
  const response = await fetch(baseUrl + url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (!response.ok) {
    const data = await response.json().catch(() => null);
    throw new SDKError("Request failed: " + response.statusText, response.status, data);
  }

  return response.json() as T;
}

`;

  for (const [resourceName, { methods }] of Object.entries(resources)) {
    const capitalizedName =
      resourceName.charAt(0).toUpperCase() + resourceName.slice(1);

    code += `
export class ${capitalizedName}Client {
  constructor(private options: ClientOptions) {}

`;

    for (const { path, httpMethod, op } of methods) {
      const methodName = getMethodName(path, httpMethod);
      const pathParams = getPathParams(op);
      const bodyType = getBodyType(op, definitions);
      const responseType = getResponseType(op, definitions);

      const paramsArg =
        pathParams.length > 0
          ? `params: { path: { ${pathParams.map((p) => `${p.name}: ${p.type}`).join("; ")} } }, `
          : "";

      const pathBase = path.replace(/\{(\w+)\}/g, "");
      const pathStr =
        pathParams.length > 0
          ? `"${pathBase}" + params.path.${pathParams[0].name}`
          : `"${path}"`;

      code += `  async ${methodName}(${paramsArg}${bodyType ? `body: ${bodyType}` : ""}): Promise<${responseType}> {
    return request<${responseType}>(${pathStr}, {
      method: "${httpMethod.toUpperCase()}"${bodyType ? ", body: JSON.stringify(body)" : ""}
    }, this.options.baseUrl);
  }

`;
    }

    code += `}

`;
  }

  code += `
export class Client {
${Object.keys(resources)
  .map(
    (name) =>
      `  public ${name}: ${name.charAt(0).toUpperCase() + name.slice(1)}Client;`,
  )
  .join("\n")}

  constructor(options: ClientOptions) {
${Object.keys(resources)
  .map(
    (name) =>
      `    this.${name} = new ${name.charAt(0).toUpperCase() + name.slice(1)}Client(options);`,
  )
  .join("\n")}
  }
}

export function createClient(options: ClientOptions): Client {
  return new Client(options);
}
`;

  return code;
}

function main() {
  const specPath = resolve(__dirname, "../../docs/swagger.json");
  const outputPath = resolve(__dirname, "../src/client.ts");

  const spec: SwaggerSpec = JSON.parse(readFileSync(specPath, "utf-8"));

  const clientCode = generateClient(spec);

  const srcDir = resolve(__dirname, "../src");
  if (!existsSync(srcDir)) {
    mkdirSync(srcDir, { recursive: true });
  }

  writeFileSync(outputPath, clientCode);
  console.log(`Generated SDK client at ${outputPath}`);
}

main();
