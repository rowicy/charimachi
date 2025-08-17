// Type declaration override for tailwind-variants
declare module "tailwind-variants/dist/config" {
  export interface TVConfig {
    [key: string]: any;
  }
}
