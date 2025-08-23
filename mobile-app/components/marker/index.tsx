import { Marker as MapMarker, type MapMarkerProps } from "react-native-maps";
import { MARKER_TYPE } from "./constants";
import { Icon } from "../ui/icon";
import { Box } from "../ui/box";
import { useMemo } from "react";

interface Props extends MapMarkerProps {
  type: keyof typeof MARKER_TYPE;
  count?: number;
  rate?: number;
}

export default function Marker({ type, count, rate, ...props }: Props) {
  const size = useMemo(() => {
    const baseSize = 22;

    if (type === "WARNING") return baseSize - 4;

    let addSize = 0;

    if (type === "VIOLATION") {
      // NOTE: 違反件数が多いもしくは違反率が高い箇所はサイズを大きくする

      if (count) addSize += count / 100;
      if (rate) addSize += (rate * 100) / 10;
    }

    return Math.round(baseSize + addSize);
  }, [type, count, rate]);

  return (
    <MapMarker {...props} pinColor={MARKER_TYPE[type].color}>
      <Box
        className="size-12 rounded-full flex flex-col justify-center items-center"
        style={{
          backgroundColor: MARKER_TYPE[type].color,
          width: size,
          height: size,
          opacity: type === "WARNING" || type === "VIOLATION" ? 0.6 : 1,
        }}
      >
        <Icon as={MARKER_TYPE[type].icon} color="white" size="xl" />
      </Box>
    </MapMarker>
  );
}
