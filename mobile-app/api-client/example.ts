import createClient from "openapi-react-query";
import createFetchClient from "openapi-fetch";
import type { paths } from "@/schema/example";
import { EXAMPLE_API_BASE_URL } from "@/env";

const fetchClient = createFetchClient<paths>({ baseUrl: EXAMPLE_API_BASE_URL });
export const $api = createClient(fetchClient);
