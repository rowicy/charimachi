import React, { useEffect, useMemo, useRef, useState } from "react";
import { View, StyleSheet, Dimensions, ActivityIndicator } from "react-native";
// @ts-ignore - react-native-maps has TypeScript compatibility issues with strict mode
import MapView, { UrlTile, Marker, Polyline } from "react-native-maps";
import { useCurrentLocation } from "@/hooks/use-location";
import { $api } from "@/api-client/api";
import { Text } from "@/components/ui/text";
import { Box } from "@/components/ui/box";
import Skeleton from "@/components/skeleton";
import {
  Checkbox,
  CheckboxGroup,
  CheckboxIcon,
  CheckboxIndicator,
  CheckboxLabel,
} from "@/components/ui/checkbox";
import { CheckIcon } from "@/components/ui/icon";

// Sample coordinates for the route
const tokyoTower = {
  latitude: 35.6586,
  longitude: 139.7454,
  latitudeDelta: 0.05,
  longitudeDelta: 0.05,
};

const convenienceStore1 = {
  latitude: 35.6762,
  longitude: 139.7654,
};

const convenienceStore2 = {
  latitude: 35.6895,
  longitude: 139.7456,
};

const tokyoSkytree = {
  latitude: 35.7101,
  longitude: 139.8107,
};

const MODES = {
  via_bike_parking: "自転車駐輪場経由",
  avoid_bus_stops: "バス停回避",
  avoid_traffic_lights: "信号回避",
} as const;

export default function MapsScreen() {
  // Use TanStack Query hook for fetching current location
  const { data: currentLocation, isLoading } = useCurrentLocation();
  const mapRef = useRef<MapView>(null);
  const [modes, setModes] = useState<string[]>([]);

  // Determine initial region - use current location if available, otherwise default to Tokyo Tower
  const initialRegion = currentLocation
    ? {
        latitude: currentLocation.latitude,
        longitude: currentLocation.longitude,
        latitudeDelta: 0.05,
        longitudeDelta: 0.05,
      }
    : tokyoTower;

  const { data: directions, isLoading: isLoadingDirections } = $api.useQuery(
    "get",
    "/directions/bicycle",
    {
      params: {
        query: {
          // NOTE: 現在地の緯度経度を使用
          start: `${currentLocation?.latitude},${currentLocation?.longitude}`,
          // NOTE: 目的地の緯度経度を使用
          end: `${tokyoSkytree.latitude},${tokyoSkytree.longitude}`,
          // NOTE: モード
          via_bike_parking: modes.includes("via_bike_parking"),
          avoid_bus_stops: modes.includes("avoid_bus_stops"),
          avoid_traffic_lights: modes.includes("avoid_traffic_lights"),
        },
      },
    },
  );

  const routeCoordinates = useMemo(() => {
    if (currentLocation) {
      return [
        {
          latitude: currentLocation.latitude,
          longitude: currentLocation.longitude,
        },
        tokyoTower,
        convenienceStore1,
        convenienceStore2,
        tokyoSkytree,
      ];
    }
    return [tokyoTower, convenienceStore1, convenienceStore2, tokyoSkytree];
  }, [currentLocation]);

  useEffect(() => {
    if (currentLocation && mapRef.current) {
      mapRef.current.animateToRegion(
        {
          latitude: currentLocation.latitude,
          longitude: currentLocation.longitude,
          latitudeDelta: 0.05,
          longitudeDelta: 0.05,
        },
        1000,
      ); // 1秒でアニメーション
    }
  }, [currentLocation]);

  return (
    <View className="flex-1 min-h-full flex items-center justify-center relative">
      {isLoading ? (
        <ActivityIndicator />
      ) : (
        <MapView ref={mapRef} style={styles.map} initialRegion={initialRegion}>
          {/* OpenStreetMap tile layer */}
          <UrlTile
            urlTemplate="https://tile.openstreetmap.org/{z}/{x}/{y}.png"
            maximumZ={19}
            minimumZ={1}
          />

          {/* Current location marker - only show if location is available */}
          {currentLocation && (
            <Marker
              coordinate={{
                latitude: currentLocation.latitude,
                longitude: currentLocation.longitude,
              }}
              title="現在地"
              description="Your current location"
              pinColor="orange"
            />
          )}

          {/* Markers for each location */}
          <Marker
            coordinate={tokyoTower}
            title="Tokyo Tower"
            description="Starting point"
            pinColor="green"
          />
          <Marker
            coordinate={convenienceStore1}
            title="Convenience Store 1"
            description="First stop"
            pinColor="blue"
          />
          <Marker
            coordinate={convenienceStore2}
            title="Convenience Store 2"
            description="Second stop"
            pinColor="blue"
          />
          <Marker
            coordinate={tokyoSkytree}
            title="Tokyo Skytree"
            description="Final destination"
            pinColor="red"
          />

          {/* Route polyline */}
          <Polyline
            coordinates={routeCoordinates}
            strokeColor="#FF6B6B"
            strokeWidth={3}
            lineCap="round"
            lineJoin="round"
          />
        </MapView>
      )}

      <Box className="z-50 absolute bottom-32 left-1/2 -translate-x-1/2 w-[90vw] p-4 bg-white rounded-lg shadow-lg">
        {/* NOTE: モード選択 */}
        <CheckboxGroup value={modes} onChange={setModes}>
          {Object.entries(MODES).map(([key, label]) => (
            <Checkbox
              key={key}
              value={key}
              isDisabled={isLoading || isLoadingDirections}
            >
              <CheckboxIndicator>
                <CheckboxIcon as={CheckIcon} />
              </CheckboxIndicator>
              <CheckboxLabel>{label}</CheckboxLabel>
            </Checkbox>
          ))}
        </CheckboxGroup>

        {/* NOTE: 距離 */}
        <SummaryItem
          label="距離"
          value={directions?.features?.[0]?.properties?.summary?.distance || 0}
          isLoading={isLoadingDirections}
        />

        {/* NOTE: 所要時間 */}
        <SummaryItem
          label="所要時間"
          // TODO: 現在地変更の度に所要時間を再計算する
          value={directions?.features?.[0]?.properties?.summary?.duration || 0}
          isLoading={isLoadingDirections}
        />
      </Box>
    </View>
  );
}

const styles = StyleSheet.create({
  map: {
    width: Dimensions.get("window").width,
    height: Dimensions.get("window").height,
  },
});

function SummaryItem({
  label,
  value,
  isLoading,
}: {
  label: string;
  value: string | number;
  isLoading: boolean;
}) {
  return (
    <Text className="color-black text-lg flex items-center">
      {label}:
      <Text className="color-black flex text-lg items-center px-2">
        {isLoading ? <Skeleton /> : value}分
      </Text>
    </Text>
  );
}
