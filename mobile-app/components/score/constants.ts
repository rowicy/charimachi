export const STATUS = {
  UNCOMFORTABLE: {
    threshold: {
      max: 49,
      min: 0,
    },
    color: "red",
    colorJp: "赤",
    textColor: "white",
  },
  CAUTION: {
    threshold: {
      max: 79,
      min: 50,
    },
    color: "yellow",
    colorJp: "黄",
    textColor: "black",
  },
  GOOD: {
    threshold: {
      max: 89,
      min: 80,
    },
    color: "green",
    colorJp: "緑",
    textColor: "white",
  },
  COMFORTABLE: {
    threshold: {
      max: 100,
      min: 90,
    },
    color: "blue",
    colorJp: "青",
    textColor: "white",
  },
} as const;
