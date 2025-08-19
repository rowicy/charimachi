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

// TODO: 目的地ダミー
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
  const { data: currentLocation, isLoading } = useCurrentLocation();
  const mapRef = useRef<MapView>(null);
  const [modes, setModes] = useState<string[]>([]);

  const initialRegion = useMemo(() => {
    if (currentLocation?.latitude && currentLocation?.longitude) {
      return {
        latitude: currentLocation.latitude,
        longitude: currentLocation.longitude,
        latitudeDelta: 0.05,
        longitudeDelta: 0.05,
      };
    }
  }, [currentLocation]);

  const { data: directions, isLoading: isLoadingDirections } = $api.useQuery(
    "get",
    "/directions/bicycle",
    {
      params: {
        query: {
          // NOTE: 現在地の緯度経度を使用
          start: `${currentLocation?.longitude},${currentLocation?.latitude}`,
          // NOTE: 目的地の緯度経度を使用
          end: `${tokyoSkytree.longitude},${tokyoSkytree.latitude}`,
          // NOTE: モード
          via_bike_parking: modes.includes("via_bike_parking"),
          avoid_bus_stops: modes.includes("avoid_bus_stops"),
          avoid_traffic_lights: modes.includes("avoid_traffic_lights"),
        },
      },
    },
    {
      enabled: !!currentLocation?.longitude && !!currentLocation?.latitude,
    },
  );

  const routeCoordinates = useMemo(() => {
    if (currentLocation && directions?.features?.[0]?.geometry?.coordinates) {
      const coordinates = directions.features[0].geometry.coordinates
        .map((coord) => ({
          latitude: coord[1],
          longitude: coord[0],
        }))
        .filter(
          (coord): coord is { latitude: number; longitude: number } =>
            typeof coord.latitude === "number" &&
            typeof coord.longitude === "number" &&
            !Number.isNaN(coord.latitude) &&
            !Number.isNaN(coord.longitude),
        );

      return [
        {
          latitude: currentLocation.latitude,
          longitude: currentLocation.longitude,
        },
        ...coordinates,
      ];
    }
    return [];
  }, [currentLocation, directions?.features]);

  const durationMinutes = useMemo(() => {
    if (directions?.features?.[0]?.properties?.summary?.duration) {
      return Math.ceil(directions.features[0].properties.summary.duration / 60);
    }

    return 0;
  }, [directions]);

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
        <MapView
          ref={mapRef}
          style={styles.map}
          // NOTE: 現在地が取得できている場合のみinitialRegionを設定
          {...(initialRegion && { initialRegion })}
        >
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

          {/* 目的地 */}
          <Marker
            coordinate={tokyoSkytree}
            title="Tokyo Skytree"
            description="Final destination"
            pinColor="red"
          />

          {/* 経路 */}
          {routeCoordinates && (
            <Polyline
              coordinates={routeCoordinates}
              strokeColor="#FF6B6B"
              strokeWidth={3}
              lineCap="round"
              lineJoin="round"
            />
          )}
        </MapView>
      )}

      <Box className="z-50 absolute bottom-32 left-1/2 -translate-x-1/2 w-[90vw] p-4 bg-white rounded-lg shadow-lg">
        {/* NOTE: モード選択 */}
        <CheckboxGroup value={modes} onChange={setModes} className="mb-4">
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
          value={
            directions?.features?.[0]?.properties?.summary?.distance ||
            "距離の取得に失敗しました"
          }
          unit="m"
          isLoading={isLoadingDirections}
        />

        {/* NOTE: 所要時間 */}
        <SummaryItem
          label="所要時間"
          value={durationMinutes}
          unit="分"
          isLoading={isLoadingDirections}
        />

        {/* NOTE: 残り時間 */}
        <SummaryItem
          label="残り時間"
          // TODO: 残り時間を算出
          value={durationMinutes}
          unit="分"
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
  unit,
  isLoading,
}: {
  label: string;
  value: string | number;
  unit: string;
  isLoading: boolean;
}) {
  return (
    <Text className="color-black text-lg flex items-center">
      {label}:&nbsp;
      <Text className="color-black flex text-lg items-center px-2">
        {isLoading ? <Skeleton /> : value}&nbsp;{unit}
      </Text>
    </Text>
  );
}
