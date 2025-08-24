import {
  Actionsheet,
  ActionsheetContent,
  ActionsheetDragIndicator,
  ActionsheetDragIndicatorWrapper,
  ActionsheetBackdrop,
} from "@/components/ui/actionsheet";
import { Button, ButtonText } from "@/components/ui/button";
import { useCallback, useMemo, useState } from "react";
import { Heading } from "../ui/heading";
import { Text } from "../ui/text";
import { STATUS } from "./constants";
import { Box } from "../ui/box";

interface Props {
  score: number;
}

export default function Score({ score }: Props) {
  const [open, setOpen] = useState(false);
  const handleClose = useCallback(() => {
    setOpen(false);
  }, []);

  const status = useMemo(() => {
    if (score < STATUS.UNCOMFORTABLE.threshold.min) {
      return STATUS.UNCOMFORTABLE;
    }
    if (score < STATUS.CAUTION.threshold.min) {
      return STATUS.CAUTION;
    }
    if (score < STATUS.GOOD.threshold.min) {
      return STATUS.GOOD;
    }
    return STATUS.COMFORTABLE;
  }, [score]);

  return (
    <>
      <Button
        onPress={() => setOpen(true)}
        className="flex flex-col w-20 h-20 absolute top-28 right-4 shadow-lg rounded-full p-2 opacity-70"
        variant="link"
        style={{
          backgroundColor: status.color,
        }}
      >
        <ButtonText className="text-md " style={{ color: status.textColor }}>
          快適度
        </ButtonText>
        <ButtonText
          className="text-3xl font-semibold leading-none -my-1 "
          style={{ color: status.textColor }}
        >
          {score}
        </ButtonText>
      </Button>

      <Actionsheet isOpen={open} onClose={handleClose}>
        <ActionsheetBackdrop />
        <ActionsheetContent>
          <ActionsheetDragIndicatorWrapper>
            <ActionsheetDragIndicator />
          </ActionsheetDragIndicatorWrapper>

          <Heading size={"xl"}>快適度とは？</Heading>

          <Text className="mt-2">
            快適度とは、表示されているルートの自転車走行における快適さを示す指標です。数値が高いほど、快適なルートです。
          </Text>

          <Heading size={"lg"} className="mt-4">
            スコアカラーについて
          </Heading>

          <Text className="mt-2">
            スコアによってカラーが変わります。
            快適なルートで、素敵なサイクリングライフを！
          </Text>

          {Object.entries(STATUS).map(([key, value]) => {
            return (
              <Box key={key} className="mt-2 ">
                <Text className="flex flex-row items-center">
                  {value.threshold.min} ~ {value.threshold.max}点:{" "}
                  <Box
                    className="size-2"
                    style={{ backgroundColor: value.color }}
                  >
                    <Text className="sr-only">{value.colorJp}</Text>
                  </Box>{" "}
                  {value.colorJp}
                </Text>
              </Box>
            );
          })}
        </ActionsheetContent>
      </Actionsheet>
    </>
  );
}
