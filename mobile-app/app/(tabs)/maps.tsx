import React from "react";
import { View, StyleSheet, Text } from "react-native";

export default function MapsScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.text}>
        Map functionality is available on mobile platforms.
      </Text>
      <Text style={styles.subtitle}>
        This feature uses OpenStreetMap with react-native-maps which requires
        iOS or Android.
      </Text>
      <Text style={styles.routeInfo}>
        Sample route: Tokyo Tower → Convenience Store 1 → Convenience Store 2 →
        Tokyo Skytree
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 20,
  },
  text: {
    fontSize: 18,
    textAlign: "center",
    marginBottom: 10,
  },
  subtitle: {
    fontSize: 14,
    textAlign: "center",
    marginBottom: 20,
    opacity: 0.7,
  },
  routeInfo: {
    fontSize: 12,
    textAlign: "center",
    fontStyle: "italic",
    opacity: 0.6,
  },
});
