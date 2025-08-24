import { Box } from "../ui/box";
import Item from "./item";
import { CREDITS } from "./constants";

export default function Credit() {
  return (
    <Box className="flex flex-col p-1 bg-white/80 text-center mt-1">
      {CREDITS.map((credit) => (
        <Item key={credit.key}>{credit}</Item>
      ))}
    </Box>
  );
}
