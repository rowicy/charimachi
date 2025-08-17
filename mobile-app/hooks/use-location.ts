import { useQuery } from "@tanstack/react-query";
import * as Location from "expo-location";

export interface LocationData {
  latitude: number;
  longitude: number;
  accuracy?: number | null;
}

export interface LocationError {
  code: string;
  message: string;
}

/**
 * Custom hook to get current location using TanStack Query
 * Handles permissions and provides error states
 */
export function useCurrentLocation() {
  return useQuery({
    queryKey: ["currentLocation"],
    queryFn: async (): Promise<LocationData> => {
      // Request location permissions
      const { status } = await Location.requestForegroundPermissionsAsync();

      if (status !== "granted") {
        throw new Error("Location permission not granted");
      }

      // Get current location
      const location = await Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.High,
      });

      return {
        latitude: location.coords.latitude,
        longitude: location.coords.longitude,
        accuracy: location.coords.accuracy,
      };
    },
    retry: false, // Don't retry permission requests
    staleTime: 1000 * 60 * 5, // Consider data fresh for 5 minutes
    gcTime: 1000 * 60 * 10, // Keep in cache for 10 minutes
  });
}
