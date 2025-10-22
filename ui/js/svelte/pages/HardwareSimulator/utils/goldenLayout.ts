import {
  ItemType,
  LayoutConfig,
  RootItemConfig,
  type GoldenLayout,
} from "golden-layout";
import { mount, type Component } from "svelte";

export function registerComponent(
  layout: GoldenLayout,
  componentName: string,
  component: Component,
) {
  layout.registerComponentFactoryFunction(componentName, (container) => {
    mount(component, { target: container.element });
  });
}

function getDefaultRootItemConfig({
  projectChipsTitle,
  editorTitle,
  simulatorTitle,
}: {
  projectChipsTitle: string;
  editorTitle: string;
  simulatorTitle: string;
}): RootItemConfig {
  return {
    type: ItemType.row,
    content: [
      {
        type: ItemType.component,
        componentType: "project-chips",
        title: projectChipsTitle,
        width: 20,
      },
      {
        type: ItemType.component,
        componentType: "editor",
        title: editorTitle,
        width: 45,
      },
      {
        type: ItemType.component,
        componentType: "simulator",
        title: simulatorTitle,
      },
    ],
  };
}

export function getLayoutConfig(options: {
  projectChipsTitle: string;
  editorTitle: string;
  simulatorTitle: string;
}): LayoutConfig {
  return {
    settings: {
      showPopoutIcon: false,
    },
    dimensions: {
      headerHeight: 36,
    },
    root: getDefaultRootItemConfig(options),
  };
}

export function disableTooltips() {
  // this removes the title attributes from all lm_tab elements
  // TODO: this only works initially, we need to observe for new lm_tabs as well
  // when tabs are moved
  document.querySelectorAll(".lm_tab").forEach((tab) => {
    tab.removeAttribute("title");
  });
}
