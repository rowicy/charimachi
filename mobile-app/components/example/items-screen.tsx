import { $api } from "@/api-client/example";
import { ActivityIndicator, FlatList, Text, View } from "react-native";
import { useNavigation } from "@react-navigation/native";
import { Screens } from "./constants";
import type { Params } from "./type";
import type { NativeStackNavigationProp } from "@react-navigation/native-stack";

type NavigationProp = NativeStackNavigationProp<
  Params,
  (typeof Screens)["ItemDetails"]
>;

export default function ItemsScreen() {
  const navigation = useNavigation<NavigationProp>();

  // Use TanStack Query hook for fetching items
  const { data: items, isLoading } = $api.useQuery("get", "/example/items");

  const handleItemPress = (id: string) => {
    navigation.navigate(Screens.ItemDetails, { id });
  };

  return (
    <View style={{ flex: 1, padding: 24 }}>
      <Text
        style={{ fontSize: 18, color: "blue" }}
        onPress={() => navigation.navigate(Screens.ApiHealth)}
      >
        API Health Check
      </Text>

      {isLoading ? (
        <ActivityIndicator />
      ) : (
        <FlatList
          data={items}
          keyExtractor={(item) => item.id.toString()}
          renderItem={({ item }) => (
            <View
              style={{
                padding: 20,
                borderBottomWidth: 1,
                borderBottomColor: "#ccc",
              }}
            >
              <Text
                style={{ fontSize: 18, color: "blue" }}
                onPress={() => handleItemPress(item.id.toString())}
              >
                {item.name}
              </Text>
            </View>
          )}
        />
      )}
    </View>
  );
}
