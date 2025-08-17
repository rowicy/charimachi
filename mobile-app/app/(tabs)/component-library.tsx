import React, { useState } from "react";
import { ScrollView, Alert } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Button, ButtonText } from "@/components/ui/button";
import { Text } from "@/components/ui/text";
import { Heading } from "@/components/ui/heading";
import { Card } from "@/components/ui/card";
import { Input, InputField } from "@/components/ui/input";

export default function ComponentLibraryScreen() {
  const [inputValue, setInputValue] = useState("");

  const handleButtonPress = () => {
    Alert.alert("Button Pressed", "This is a Gluestack-UI button!");
  };

  const handleInputChange = (text: string) => {
    setInputValue(text);
  };

  return (
    <SafeAreaView className="flex-1 bg-background-0">
      <ScrollView className="flex-1 p-4">
        <Heading size="2xl" className="mb-6 text-center text-typography-900">
          Gluestack-UI Components
        </Heading>

        <Text className="mb-6 text-center text-typography-700">
          This screen showcases various Gluestack-UI components with Tailwind
          CSS styling.
        </Text>

        {/* Typography Section */}
        <Card className="mb-6 p-4">
          <Heading size="lg" className="mb-4 text-typography-900">
            Typography
          </Heading>

          <Heading size="md" className="mb-2">
            This is a heading
          </Heading>

          <Text size="lg" className="mb-2">
            Large text with custom styling
          </Text>
          <Text size="md" className="mb-2">
            Medium text (default size)
          </Text>
          <Text size="sm" className="mb-2 text-typography-600">
            Small text in muted color
          </Text>

          <Text bold className="mb-2">
            Bold text example
          </Text>
          <Text italic className="mb-2">
            Italic text example
          </Text>
          <Text underline className="mb-2">
            Underlined text example
          </Text>
        </Card>

        {/* Buttons Section */}
        <Card className="mb-6 p-4">
          <Heading size="lg" className="mb-4 text-typography-900">
            Buttons
          </Heading>

          <Button
            onPress={handleButtonPress}
            className="mb-3"
            action="primary"
            variant="solid"
            size="md"
          >
            <ButtonText>Primary Button</ButtonText>
          </Button>

          <Button
            onPress={handleButtonPress}
            className="mb-3"
            action="secondary"
            variant="solid"
            size="md"
          >
            <ButtonText>Secondary Button</ButtonText>
          </Button>

          <Button
            onPress={handleButtonPress}
            className="mb-3"
            action="positive"
            variant="outline"
            size="md"
          >
            <ButtonText>Success Outline</ButtonText>
          </Button>

          <Button
            onPress={handleButtonPress}
            className="mb-3"
            action="negative"
            variant="outline"
            size="md"
          >
            <ButtonText>Error Outline</ButtonText>
          </Button>

          <Button
            onPress={handleButtonPress}
            action="primary"
            variant="link"
            size="md"
          >
            <ButtonText>Link Button</ButtonText>
          </Button>
        </Card>

        {/* Input Section */}
        <Card className="mb-6 p-4">
          <Heading size="lg" className="mb-4 text-typography-900">
            Input Field
          </Heading>

          <Input className="mb-4">
            <InputField
              placeholder="Enter some text..."
              value={inputValue}
              onChangeText={handleInputChange}
            />
          </Input>

          <Text size="sm" className="text-typography-600">
            Current value: {inputValue || "No text entered"}
          </Text>
        </Card>

        {/* Color Palette Section */}
        <Card className="mb-6 p-4">
          <Heading size="lg" className="mb-4 text-typography-900">
            Color System
          </Heading>

          <Text className="mb-3 text-primary-500">Primary color text</Text>
          <Text className="mb-3 text-secondary-500">Secondary color text</Text>
          <Text className="mb-3 text-success-500">Success color text</Text>
          <Text className="mb-3 text-error-500">Error color text</Text>
          <Text className="mb-3 text-warning-500">Warning color text</Text>
          <Text className="mb-3 text-info-500">Info color text</Text>
        </Card>

        {/* Tailwind Classes Section */}
        <Card className="mb-6 p-4">
          <Heading size="lg" className="mb-4 text-typography-900">
            Tailwind CSS Classes
          </Heading>

          <Text className="mb-2 text-center text-lg font-bold text-primary-600">
            Centered bold primary text
          </Text>

          <Text className="mb-2 italic text-secondary-700">
            Italic secondary text
          </Text>

          <Text className="rounded-lg bg-success-100 p-3 text-success-800">
            Text with background and padding
          </Text>
        </Card>
      </ScrollView>
    </SafeAreaView>
  );
}
