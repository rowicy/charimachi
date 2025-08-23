import { Marker as MapMarker, type MapMarkerProps } from "react-native-maps";
import { MARKER_TYPE } from "./constants";
import { Icon } from "../ui/icon";
import { Box } from "../ui/box";

interface Props extends MapMarkerProps {
  type: keyof typeof MARKER_TYPE;
}

export default function Marker({ type, ...props }: Props) {
  const size = type === "WARNING" ? 40 : 48;

  return (
    <MapMarker {...props} pinColor={MARKER_TYPE[type].color}>
      <Box
        className="size-12 rounded-full flex flex-col justify-center items-center"
        style={{
          backgroundColor: MARKER_TYPE[type].color,
          width: size,
          height: size,
        }}
      >
        <Icon as={MARKER_TYPE[type].icon} color="white" size="xl" />
      </Box>
    </MapMarker>
  );
}
