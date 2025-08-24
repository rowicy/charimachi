import { Box } from "../ui/box";
import { Text } from "../ui/text";

interface Props {
  width?: number;
  height?: number;
}

export default function Skeleton({ width = 12, height = 12 }: Props) {
  return (
    <Box className="px-1">
      <Box className="animate-pulse bg-gray-300" style={{ width, height }}>
        <Text className="sr-only">Loading...</Text>
      </Box>
    </Box>
  );
}
