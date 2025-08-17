import { $api } from "@/api-client/api";
import { Text } from "../ui/text";
import { Heading } from "../ui/heading";
import { View } from "react-native";

export default function ApiHealthScreen() {
  const { data, isLoading } = $api.useQuery("get", "/health");

  return (
    <View style={{ flex: 1, padding: 24 }}>
      <Heading>API Health Status</Heading>

      {isLoading ? (
        <Text>Loading...</Text>
      ) : (
        <Text style={{ fontSize: 16, color: "green" }}>
          {JSON.stringify(data, null, 2)}
        </Text>
      )}
    </View>
  );
}
