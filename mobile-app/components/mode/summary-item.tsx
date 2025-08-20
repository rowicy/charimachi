import { Text } from "@/components/ui/text";
import Skeleton from "@/components/skeleton";

interface Props {
  label: string;
  value?: number;
  unit: string;
  loading: boolean;
  error?: boolean;
}

export default function SummaryItem({
  label,
  value,
  unit,
  loading,
  error,
}: Props) {
  return (
    <Text className="color-black text-lg flex items-center">
      {label}:&nbsp;
      <Text className="color-black flex text-lg items-center px-2">
        {loading && <Skeleton />}
        {value && `${value} ${unit}`}
        {error && "取得に失敗しました"}
      </Text>
    </Text>
  );
}
