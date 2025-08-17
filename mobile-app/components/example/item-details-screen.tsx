import { $api } from "@/api-client/example";
import { ActivityIndicator, Text, View } from "react-native";
import { useNavigation } from "@react-navigation/native";
import { Screens } from "./constants";
import type { Params } from "./type";
import type {
  NativeStackNavigationProp,
  NativeStackScreenProps,
} from "@react-navigation/native-stack";

type NavigationProp = NativeStackNavigationProp<
  Params,
  (typeof Screens)["Items"]
>;

type ItemDetailsScreenProps = NativeStackScreenProps<
  Params,
  (typeof Screens)["ItemDetails"]
>;

export default function ItemDetailsScreen({ route }: ItemDetailsScreenProps) {
  const navigation = useNavigation<NavigationProp>();

  // Use TanStack Query hook for fetching item by ID
  const { data: item, isLoading } = $api.useQuery(
    "get",
    "/example/items/{id}",
    {
      params: { path: { id: route.params?.id || "" } },
    },
    {
      enabled: !!route.params?.id, // Only run query if ID is available
    },
  );

  return (
    <View>
      {isLoading ? (
        <ActivityIndicator />
      ) : (
        <View style={{ padding: 20 }}>
          <View style={{ marginBottom: 10 }}>
            <Text style={{ color: "red", fontSize: 18, fontWeight: "bold" }}>
              ID:
            </Text>
            <Text style={{ color: "red", fontSize: 18, marginLeft: 10 }}>
              {item?.id}
            </Text>
          </View>
          <View style={{ marginBottom: 10 }}>
            <Text style={{ color: "red", fontSize: 18, fontWeight: "bold" }}>
              Name:
            </Text>
            <Text style={{ color: "red", fontSize: 18, marginLeft: 10 }}>
              {item?.name}
            </Text>
          </View>
          <View style={{ marginBottom: 10 }}>
            <Text style={{ color: "red", fontSize: 18, fontWeight: "bold" }}>
              Price:
            </Text>
            <Text style={{ color: "red", fontSize: 18, marginLeft: 10 }}>
              {item?.price}
            </Text>
          </View>
        </View>
      )}

      <Text
        style={{
          fontSize: 18,
          color: "blue",
          textAlign: "center",
          marginTop: 20,
        }}
        onPress={() => navigation.navigate(Screens.Items)}
      >
        Back to Items List
      </Text>
    </View>
  );
}
