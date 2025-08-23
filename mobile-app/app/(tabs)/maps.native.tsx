import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { StyleSheet, Dimensions, ActivityIndicator } from "react-native";
// @ts-ignore - react-native-maps has TypeScript compatibility issues with strict mode
import MapView, { UrlTile, Marker, Polyline } from "react-native-maps";
import { useCurrentLocation } from "@/hooks/use-location";
import { $api } from "@/api-client/api";
import { Box } from "@/components/ui/box";
import type { components } from "@/schema/api";
import Search from "@/components/search";
import Mode from "@/components/mode";
import { Text } from "@/components/ui/text";
import { Link, LinkText } from "@/components/ui/link";
import Score from "@/components/score";

export default function MapsScreen() {
  const { data: currentLocation, isLoading } = useCurrentLocation();
  const mapRef = useRef<MapView>(null);
  const [openSearch, setOpenSearch] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [debouncedKeyword, setDebouncedKeyword] = useState(keyword);
  const [destination, setDestination] = useState<
    components["schemas"]["main.SearchResponse"] | null
  >(null);
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

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedKeyword(keyword);
      setOpenSearch(!!keyword && keyword.length > 0);
    }, 1000);
    return () => clearTimeout(handler);
  }, [keyword]);

  const { data: destinations, isLoading: isLoadingDestinations } =
    $api.useQuery(
      "get",
      "/search",
      {
        params: {
          query: {
            q: debouncedKeyword,
          },
        },
      },
      {
        enabled: !!debouncedKeyword && debouncedKeyword.length > 0,
      },
    );

  const handleDestinationSelect = useCallback(
    (destination: components["schemas"]["main.DirectionsResponse"]) => {
      setDestination(destination);
      setOpenSearch(false);
    },
    [],
  );

  const {
    data: directions,
    isLoading: isLoadingDirections,
    isError: isErrorDirections,
  } = $api.useQuery(
    "get",
    "/directions/bicycle",
    {
      params: {
        query: {
          // NOTE: 現在地の緯度経度を使用
          start: `${currentLocation?.longitude},${currentLocation?.latitude}`,
          // NOTE: 目的地の緯度経度を使用
          end: `${destination?.lon},${destination?.lat}`,
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

  const distance = useMemo(() => {
    // NOTE: 1000m以上の場合はkmで返す
    if (directions?.features?.[0]?.properties?.summary?.distance) {
      if (directions.features[0].properties.summary.distance >= 1000) {
        return {
          value: Number(
            (directions.features[0].properties.summary.distance / 1000).toFixed(
              1,
            ),
          ),
          unit: "km",
        };
      }
      // NOTE: 1000m未満の場合はmで返す
      return {
        value: Math.round(directions.features[0].properties.summary.distance),
        unit: "m",
      };
    }
  }, [directions?.features]);

  const duration = useMemo(() => {
    if (directions?.features?.[0]?.properties?.summary?.duration) {
      // NOTE: 60分(3600秒)未満の場合は分で返す
      if (directions.features[0].properties.summary.duration < 3600) {
        return {
          value: Math.ceil(
            directions.features[0].properties.summary.duration / 60,
          ),
          unit: "分",
        };
      }

      // NOTE: n.n時間
      return {
        value: Math.ceil(
          directions.features[0].properties.summary.duration / 3600,
        ),
        unit: "時間",
      };
    }

    return {
      value: 0,
      unit: "分",
    };
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

  // NOTE: 目的地が設定されたら拡大位置をbbox基準に変更
  useEffect(() => {
    if (destination && currentLocation && mapRef.current) {
      const destLat = Number(destination.lat);
      const destLon = Number(destination.lon);
      const currLat = currentLocation.latitude;
      const currLon = currentLocation.longitude;

      // 両点を含む範囲を計算
      const minLat = Math.min(destLat, currLat);
      const maxLat = Math.max(destLat, currLat);
      const minLon = Math.min(destLon, currLon);
      const maxLon = Math.max(destLon, currLon);

      // 中心点とデルタを計算
      const centerLat = (minLat + maxLat) / 2;
      const centerLon = (minLon + maxLon) / 2;
      const latDelta = Math.max((maxLat - minLat) * 1.5, 0.005);
      const lonDelta = Math.max((maxLon - minLon) * 1.5, 0.005);

      mapRef.current.animateToRegion(
        {
          latitude: centerLat,
          longitude: centerLon,
          latitudeDelta: latDelta,
          longitudeDelta: lonDelta,
        },
        1000,
      );
    }
  }, [destination, currentLocation]);

  return (
    <Box className="flex-1 min-h-full flex items-center justify-center relative">
      {isLoading ? (
        <ActivityIndicator />
      ) : (
        <>
          {/* NOTE: 目的地検索 */}
          <Search
            keyword={keyword}
            setKeyword={setKeyword}
            open={openSearch}
            loading={isLoadingDestinations}
            destinations={destinations}
            onSelect={handleDestinationSelect}
          />

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
            {destination?.display_name &&
              destination?.lat &&
              destination?.lon && (
                <Marker
                  coordinate={{
                    latitude: Number(destination?.lat),
                    longitude: Number(destination?.lon),
                  }}
                  title={destination?.display_name}
                  description="目的地"
                  pinColor="red"
                />
              )}

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

          {/* NOTE: スコア */}
          {directions?.comfort_score && (
            <Score score={directions?.comfort_score} />
          )}

          <Box className="absolute bottom-28 left-1/2 -translate-x-1/2 w-[90vw] flex items-end flex-col">
            {/* NOTE: モード選択 */}
            <Mode
              loading={isLoading || isLoadingDirections}
              distance={distance}
              duration={duration}
              modes={modes}
              setModes={setModes}
              error={!!destination && isErrorDirections}
              destination={!!destination}
            />

            <Box className="flex flex-row items-center justify-center p-1 text-gray-500 bg-white/80 text-center mt-1">
              <Text className="text-sm">&copy;&nbsp;</Text>
              <Link href="https://www.openstreetmap.org/copyright">
                <LinkText size="sm">OpenStreetMap contributors</LinkText>
              </Link>
            </Box>
          </Box>
        </>
      )}
    </Box>
  );
}

const styles = StyleSheet.create({
  map: {
    width: Dimensions.get("window").width,
    height: Dimensions.get("window").height,
  },
});
