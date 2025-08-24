import type { ReactNode } from "react";
import { Box } from "../ui/box";

interface Props {
  children: ReactNode;
}

export default function Item({ children }: Props) {
  return <Box className="flex flex-row items-center">{children}</Box>;
}
