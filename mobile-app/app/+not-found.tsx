import { Stack } from "expo-router";
import { View } from "react-native";
import { Heading } from "@/components/ui/heading";

export default function NotFoundScreen() {
  return (
    <>
      <Stack.Screen options={{ title: "Oops!" }} />
      <View>
        <Heading>This screen does not exist.</Heading>
      </View>
    </>
  );
}
