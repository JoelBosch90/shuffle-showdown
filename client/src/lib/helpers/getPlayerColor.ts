import { textToNumber } from './textToNumber';

const PLAYER_COLOURS = [
  'folly', 'satin-sheen-gold', 'celestial-blue', 'persian-green', 'verdigris', 'tropical-indigo', 'claret', 'princeton-orange', 'emerald', 'indigo-dye'
];

export const getPlayerColor = (playerId: string): string => {
  const lastPart = playerId.split('-').at(-1) ?? '';

  return PLAYER_COLOURS[textToNumber(lastPart) % PLAYER_COLOURS.length];
};