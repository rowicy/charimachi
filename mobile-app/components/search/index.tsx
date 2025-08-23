import type { components } from "@/schema/api";
import Skeleton from "../skeleton";
import { Box } from "../ui/box";
import { Button } from "../ui/button";
import { FormControl } from "../ui/form-control";
import { SearchIcon } from "../ui/icon";
import { Input, InputField, InputIcon, InputSlot } from "../ui/input";
import { Text } from "../ui/text";
import Item from "./item";

export interface SearchProps {
  keyword: string;
  setKeyword: (keyword: string) => void;
  open: boolean;
  loading: boolean;
  destinations?: components["schemas"]["main.SearchResponse"][];
  onSelect: (destination: components["schemas"]["main.SearchResponse"]) => void;
}

export default function Search({
  keyword,
  setKeyword,
  open,
  loading,
  destinations,
  onSelect,
}: SearchProps) {
  return (
    <Box className="absolute top-16 left-1/2 -translate-x-1/2 w-[90vw] z-50">
      {/* NOTE: 入力欄 */}
      <FormControl className="shadow-lg">
        <Input className="bg-white outline-none border-white">
          <InputSlot className="pl-3">
            <InputIcon as={SearchIcon} />
          </InputSlot>
          <InputField
            placeholder="目的地を入力"
            onChangeText={(text) => setKeyword(text)}
            value={keyword}
            className="text-black"
          />
        </Input>
      </FormControl>

      {/* NOTE: サジェスト */}
      {open && (
        <Box className="bg-white rounded-b-lg shadow-lg mt-1">
          {loading ? (
            // skeletonループで3つ表示
            <>
              {Array.from({ length: 3 }, (_, index) => (
                <Item key={index}>
                  <Skeleton width={200} height={20} />
                </Item>
              ))}
            </>
          ) : Array.isArray(destinations) && destinations.length > 0 ? (
            destinations.map((destination) => {
              if (!destination.display_name) return;

              return (
                <Item key={destination.place_id}>
                  <Button
                    className="h-auto justify-start"
                    variant="link"
                    onPress={() => onSelect(destination)}
                  >
                    <Text className="text-black">
                      {destination.display_name}
                    </Text>
                  </Button>
                </Item>
              );
            })
          ) : (
            <Text className="p-2 text-black border-gray-200">
              見つかりませんでした。
            </Text>
          )}
        </Box>
      )}
    </Box>
  );
}
