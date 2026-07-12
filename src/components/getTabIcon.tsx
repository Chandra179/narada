import type { SVGProps } from 'react';
import {
  RomanceIcon,
  HealthIcon,
  CareerIcon,
  WealthIcon,
  StarIcon,
} from './icons';

const ICON_MAP: Record<string, React.ElementType<SVGProps<SVGSVGElement>>> = {
  heart: RomanceIcon,
  person: HealthIcon,
  briefcase: CareerIcon,
  circles: WealthIcon,
  star: StarIcon,
};

export function getTabIcon(iconKey: string): React.ElementType<SVGProps<SVGSVGElement>> {
  return ICON_MAP[iconKey] ?? StarIcon;
}
