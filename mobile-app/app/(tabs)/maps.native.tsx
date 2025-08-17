import React, { useEffect, useMemo, useRef } from "react";
import { View, StyleSheet, Dimensions, ActivityIndicator } from "react-native";
// @ts-ignore - react-native-maps has TypeScript compatibility issues with strict mode
import MapView, { UrlTile, Marker, Polyline } from "react-native-maps";
import { useCurrentLocation } from "@/hooks/use-location";

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

export default function MapsScreen() {
  // Use TanStack Query hook for fetching current location
  const { data: currentLocation, isLoading } = useCurrentLocation();
  const mapRef = useRef<MapView>(null);

  // Determine initial region - use current location if available, otherwise default to Tokyo Tower
  const initialRegion = currentLocation
    ? {
        latitude: currentLocation.latitude,
        longitude: currentLocation.longitude,
        latitudeDelta: 0.05,
        longitudeDelta: 0.05,
      }
    : tokyoTower;

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
    <View className="flex-1 min-h-full flex items-center justify-center">
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
    </View>
  );
}

const styles = StyleSheet.create({
  map: {
    width: Dimensions.get("window").width,
    height: Dimensions.get("window").height,
  },
});
