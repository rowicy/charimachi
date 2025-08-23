import {
  Actionsheet,
  ActionsheetContent,
  ActionsheetDragIndicator,
  ActionsheetDragIndicatorWrapper,
  ActionsheetBackdrop,
} from "@/components/ui/actionsheet";
import {
  Checkbox,
  CheckboxGroup,
  CheckboxIcon,
  CheckboxIndicator,
  CheckboxLabel,
} from "@/components/ui/checkbox";
import { MODES } from "@/constants/modes";
import { CheckIcon } from "../ui/icon";
import { useCallback, useState } from "react";
import { Fab, FabLabel, FabIcon } from "@/components/ui/fab";
import { AddIcon } from "@/components/ui/icon";

interface Props {
  loading: boolean;
  modes: string[];
  setModes: (modes: string[]) => void;
}

export default function Mode({ loading, modes, setModes }: Props) {
  const [open, setOpen] = useState(false);
  const handleClose = useCallback(() => {
    setOpen(false);
  }, []);

  return (
    <>
      <Fab
        size="md"
        placement="bottom right"
        onPress={() => setOpen(true)}
        className="bottom-28"
      >
        <FabIcon as={AddIcon} />
        <FabLabel>モード切替</FabLabel>
      </Fab>

      <Actionsheet isOpen={open} onClose={handleClose}>
        <ActionsheetBackdrop />
        <ActionsheetContent>
          <ActionsheetDragIndicatorWrapper>
            <ActionsheetDragIndicator />
          </ActionsheetDragIndicatorWrapper>

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
        </ActionsheetContent>
      </Actionsheet>
    </>
  );
}
