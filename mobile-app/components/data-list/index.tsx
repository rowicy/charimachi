import React from "react";
import { Link, LinkText } from "../ui/link";
import { DATA_LIST } from "./constants";

export default function DataList() {
  return (
    <>
      {DATA_LIST.map((data) => (
        <Link href={data.link} key={data.link} className="mt-2">
          <LinkText>{data.label}</LinkText>
        </Link>
      ))}
    </>
  );
}
