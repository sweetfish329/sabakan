import { applicationConfig, type Meta, type StoryObj } from "@storybook/angular";
import { provideRouter } from "@angular/router";
import { provideHttpClient } from "@angular/common/http";
import { provideAnimationsAsync } from "@angular/platform-browser/animations/async";

import { ModListComponent } from "./mod-list";

const meta: Meta<ModListComponent> = {
  title: "Features/Mods/ModList",
  component: ModListComponent,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
  },
  decorators: [
    applicationConfig({
      providers: [provideRouter([]), provideHttpClient(), provideAnimationsAsync()],
    }),
  ],
};

export default meta;
type Story = StoryObj<ModListComponent>;

/**
 * Default mod list showing available mods.
 */
export const Default: Story = {};
