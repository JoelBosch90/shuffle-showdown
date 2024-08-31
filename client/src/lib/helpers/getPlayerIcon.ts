import { textToNumber } from './textToNumber';

const PLAYER_ICONS = [
  'hippo', 'fish', 'dragon', 'otter', 'worm', 'spider', 'shrimp', 'mosquito', 'locust', 'horse', 'frog', 'dove', 'dog', 'crow', 'cow', 'cat',
  'car', 'plane', 'truck', 'bicycle', 'rocket', 'ship', 'motorcycle', 'train', 'tractor', 'taxi', 'sleigh', 'bus'
];

export const getPlayerIcon = (playerId: string): string => {
  const secondLastPart = playerId.split('-').at(-2) ?? '';

  return PLAYER_ICONS[textToNumber(secondLastPart) % PLAYER_ICONS.length];
};