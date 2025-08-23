interface Data {
  label: string;
  link: string;
}

const openDataSiteName = "TOKYO OPEN DATA";

export const DATA_LIST: Data[] = [
  {
    label: "OpenStreetMap",
    link: "https://www.openstreetmap.org/",
  },
  {
    label: `${openDataSiteName} - 重点取締場所一覧 交通違反重点取締場所一覧表`,
    link: "https://spec.api.metro.tokyo.lg.jp/spec/t000022d1700000024-29a128f7bb366ba2832927fac7feeaa4-0?lang=ja",
  },
  {
    label: `${openDataSiteName} - 交通量統計表`,
    link: "https://catalog.data.metro.tokyo.lg.jp/dataset/t000022d0000000035",
  },
];
