import type { ReactNode } from "react";
import { Box } from "../ui/box";

interface Props {
  children: ReactNode;
}

export default function Item({ children }: Props) {
  return (
    <Box className="py-2 px-4 [&:not(:first-child)]:border-b border-gray-200">
      {children}
    </Box>
  );
}
