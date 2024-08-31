export const textToNumber = (text: string): number => {
  return [...text].reduce((accumulator, character) => accumulator + character.charCodeAt(0), 0) ?? 0;
};