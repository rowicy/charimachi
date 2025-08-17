import { createNativeStackNavigator } from "@react-navigation/native-stack";
import ItemsScreen from "./items-screen";
import ItemDetailsScreen from "./item-details-screen";
import ApiHealthScreen from "./api-health-screen";
import { Screens } from "./constants";
import type { Params } from "./type";

const Stack = createNativeStackNavigator<Params>();

export default function Example() {
  return (
    <Stack.Navigator id={undefined} initialRouteName={Screens.Items}>
      <Stack.Screen name={Screens.Items} component={ItemsScreen} />
      <Stack.Screen name={Screens.ItemDetails} component={ItemDetailsScreen} />
      <Stack.Screen name={Screens.ApiHealth} component={ApiHealthScreen} />
    </Stack.Navigator>
  );
}
