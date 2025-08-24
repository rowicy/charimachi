import { Link, LinkText } from "../ui/link";
import { Text } from "../ui/text";
import { Box } from "../ui/box";

const SIZE = "xs";

export const CREDITS = [
  <>
    <Text size={SIZE} className="text-gray-500" key={"osm-copyright"}>
      &copy;&nbsp;
    </Text>
    <Link href="https://www.openstreetmap.org/copyright" key={"osm-link"}>
      <LinkText size={SIZE}>OpenStreetMap contributors</LinkText>
    </Link>
  </>,
  <Box
    key={"bus-stop-information"}
    className="flex flex-row flex-wrap border-t-gray-300 border-t pt-1 mt-1"
  >
    <Text size={SIZE} className="text-gray-500 w-full">
      このアプリケーションは、以下の著作物を改変して利用しています。
    </Text>
    <Text size={SIZE} className="text-gray-500">
      東京都交通局・公共交通オープンデータ協議会、バス停情報 / Bus stop
      information、
    </Text>
    <Link href="https://creativecommons.org/licenses/by/4.0/deed.ja">
      <LinkText size={SIZE}>
        クリエイティブ・コモンズ・ライセンス　表示4.0国際
      </LinkText>
    </Link>
  </Box>,
];
