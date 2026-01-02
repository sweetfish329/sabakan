import globals from "globals";
import pluginJs from "@eslint/js";
import tseslint from "typescript-eslint";
import pluginOxlint from "eslint-plugin-oxlint";
import { FlatCompat } from "@eslint/eslintrc";
import path from "path";
import { fileURLToPath } from "url";
import pluginImport from "eslint-plugin-import";
import pluginJsdoc from "eslint-plugin-jsdoc";
import pluginUnicorn from "eslint-plugin-unicorn";
import pluginDeprecation from "eslint-plugin-deprecation";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

export default tseslint.config(
  { files: ["**/*.{js,mjs,cjs,ts}"] },
  { languageOptions: { globals: globals.browser } },
  {
    plugins: {
      import: pluginImport,
      jsdoc: pluginJsdoc,
      unicorn: pluginUnicorn,
      deprecation: pluginDeprecation,
    },
  },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  ...compat.extends("plugin:storybook/recommended"),
  ...pluginOxlint.configs["flat/recommended"],
  {
    ignores: ["dist/**", "coverage/**", "storybook-static/**", ".angular/**"],
  },
);
