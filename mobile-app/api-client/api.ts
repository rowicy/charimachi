import createClient from "openapi-react-query";
import createFetchClient from "openapi-fetch";
import type { paths } from "@/schema/api";
import { API_BASE_URL } from "@/env";

const fetchClient = createFetchClient<paths>({ baseUrl: API_BASE_URL });
export const $api = createClient(fetchClient);
