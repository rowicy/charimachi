import { Box } from "@/components/ui/box";
import {
  Checkbox,
  CheckboxGroup,
  CheckboxIcon,
  CheckboxIndicator,
  CheckboxLabel,
} from "@/components/ui/checkbox";
import { MODES } from "@/constants/modes";
import SummaryItem from "./summary-item";
import { CheckIcon } from "../ui/icon";

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
  modes: string[];
  setModes: (modes: string[]) => void;
  error?: boolean;
  destination: boolean;
}

export default function Mode({
  loading,
  distance,
  duration,
  modes,
  setModes,
  error,
  destination,
}: Props) {
  return (
    <Box className="z-50 p-4 bg-white/80 rounded-lg shadow-lg w-full">
      {/* NOTE: モード選択 */}
      <CheckboxGroup value={modes} onChange={setModes}>
        {Object.entries(MODES).map(([key, label]) => (
          <Checkbox
            key={key}
            value={key}
            isDisabled={loading}
            className="bg-transparent"
          >
            <CheckboxIndicator>
              <CheckboxIcon as={CheckIcon} />
            </CheckboxIndicator>
            {/* 強制的にテキスト色を黒に固定 */}
            <CheckboxLabel style={{ color: "#000" }}>{label}</CheckboxLabel>
          </Checkbox>
        ))}
      </CheckboxGroup>

      {destination && (
        <Box className="mt-4">
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
      )}
    </Box>
  );
}
