import { Box } from "@/components/ui/box";
import SummaryItem from "./summary-item";

interface Props {
  loading: boolean;
  distance?: {
    value: number;
    unit: string;
  };
  duration?: {
    value: number;
    unit: string;
  };
  error?: boolean;
  destination: boolean;
}

export default function Result({
  loading,
  distance,
  duration,
  error,
  destination,
}: Props) {
  if (!destination) return null;

  return (
    <Box className="z-50 p-4 bg-white/80 rounded-lg shadow-lg w-full">
      {/* NOTE: 距離 */}
      <SummaryItem
        label="距離"
        value={distance?.value}
        unit={distance?.unit}
        loading={loading}
      />

      {/* NOTE: 所要時間 */}
      <SummaryItem
        label="所要時間"
        value={duration?.value}
        unit={duration?.unit}
        loading={loading}
        error={error}
      />

      {/* NOTE: 残り時間 */}
      <SummaryItem
        label="残り時間"
        // TODO: 残り時間を算出
        value={duration?.value}
        unit={duration?.unit}
        loading={loading}
        error={error}
      />
    </Box>
  );
}
