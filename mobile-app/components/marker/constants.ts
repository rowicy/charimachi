import { PlayIcon, CheckCircleIcon, AlertCircleIcon } from "../ui/icon";

export const MARKER_TYPE = {
  START: {
    label: "現在地",
    color: "green",
    icon: PlayIcon,
  },
  GOAL: {
    label: "目的地",
    color: "blue",
    icon: CheckCircleIcon,
  },
  WARNING: {
    label: "警告",
    color: "red",
    icon: AlertCircleIcon,
  },
};
