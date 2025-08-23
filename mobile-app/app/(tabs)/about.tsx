import { Heading } from "@/components/ui/heading";
import { Text } from "@/components/ui/text";
import { SafeAreaView, ScrollView } from "react-native";

export default function AboutScreen() {
  return (
    <SafeAreaView className="flex-1 bg-background-0">
      <ScrollView className="flex-1 p-4">
        <Heading size="2xl" className="mb-6 text-center text-typography-900">
          About
        </Heading>
        <Text size="md">foo bar</Text>
      </ScrollView>
    </SafeAreaView>
  );
}
