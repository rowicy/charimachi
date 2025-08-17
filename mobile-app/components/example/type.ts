import { Screens } from "./constants";

export type Params = {
  [Screens.Items]: undefined;
  [Screens.ItemDetails]: { id: string };
  [Screens.ApiHealth]: undefined;
};
