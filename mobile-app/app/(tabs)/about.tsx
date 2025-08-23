import { Heading } from "@/components/ui/heading";
import { Text } from "@/components/ui/text";
import type { ReactNode } from "react";
import { SafeAreaView, ScrollView } from "react-native";
import { Image } from "@/components/ui/image";
import { Center } from "@/components/ui/center";
import DataList from "@/components/data-list";
import { Link, LinkText } from "@/components/ui/link";

export default function AboutScreen() {
  return (
    <SafeAreaView className="flex-1 bg-background-0">
      <ScrollView className="flex-1 p-4">
        <Heading size="2xl" className="mb-6 text-center text-typography-900">
          ChariMachi
        </Heading>

        <Center>
          <Image
            size="lg"
            source={require("@/assets/images/icon.png")}
            alt="ChariMachiのロゴ"
          />
        </Center>

        <Description>
          ChariMachiは、自転車の車道通行を目的とした適切なルート提案アプリです。
        </Description>
        <Description>
          このアプリは、ユーザーが自転車を安全かつ快適に利用できるよう、最適なルートを提案します。
        </Description>

        <Title>利用しているデータ</Title>
        <DataList />

        <Center className="mt-4">
          <Title>Produced by</Title>

          <Link href="https://www.rowicy.com/">
            <LinkText className="text-center" size="xl">
              Rowicy
            </LinkText>

            <Image
              size="lg"
              source={require("@/assets/images/rowicy.png")}
              alt="rowicyのロゴ"
              className="mt-2"
            />
          </Link>
        </Center>
      </ScrollView>
    </SafeAreaView>
  );
}

function Title({ children }: { children: ReactNode }) {
  return (
    <Text size="lg" className="mt-4 font-bold">
      {children}
    </Text>
  );
}

function Description({ children }: { children: ReactNode }) {
  return (
    <Text size="md" className="mt-4">
      {children}
    </Text>
  );
}
