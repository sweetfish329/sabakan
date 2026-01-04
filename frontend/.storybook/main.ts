import type { StorybookConfig } from "@storybook/angular";

const config: StorybookConfig = {
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|mjs|ts|tsx)"],
  staticDirs: [{ from: "../src/assets", to: "/assets" }],
  addons: [
    "@storybook/addon-links",
    "@storybook/addon-essentials",

    "@storybook/addon-a11y",
    "@storybook/experimental-addon-test",
  ],
  framework: {
    name: "@storybook/angular",
    options: {},
  },
  /* docs: {
    autodocs: "tag",
  }, */
};
export default config;
