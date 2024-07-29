import fs from 'fs';

// Enums
export enum Race {
  Caucasian = 'Caucasian',
  Mongoloid = 'Mongoloid',
  Negroid = 'Negroid',
  Android = 'Android',
}

export enum Emotion {
  Neutral = 'Neutral',
  Happiness = 'Happiness',
  Anger = 'Anger',
  Surprise = 'Surprise',
  Fear = 'Fear',
  Sadness = 'Sadness',
  Disgust = 'Disgust',
}

export enum Age {
  Baby = 'Baby',
  Kid = 'Kid',
  Teenager = 'Teenager',
  YoungAdult = '20-30s',
  MaturedAdult = '40-50s',
  Senior = 'Senior',
}

export enum Gender {
  Male = 'Male',
  Female = 'Female',
}

// Superhero type definition
export type Superhero = {
  file_name: string;
  bbox: number[];
  gender: Gender;
  emotion: Emotion;
  age: Age;
  race: Race;
};

// Helper function to get random enum value
export function getRandomEnumValue<T>(enumObj: T): T[keyof T] {
  const enumValues = Object.keys(enumObj as any) as Array<keyof T>;
  const randomIndex = Math.floor(Math.random() * enumValues.length);
  return enumObj[enumValues[randomIndex]];
}

export function readSuperheroesFromFile(filePath: string): Superhero[] {
  const data = fs.readFileSync(filePath, 'utf-8');
  return JSON.parse(data) as Superhero[];
}

// Helper functions to check the criteria
export function isBalanceGuardians(group: Superhero[]): boolean {
  if (group.length !== 4) return false;
  const maleCount = group.filter((hero) => hero.gender === Gender.Male).length;
  const femaleCount = group.filter(
    (hero) => hero.gender === Gender.Female
  ).length;
  const uniqueRaces = new Set(group.map((hero) => hero.race));

  return maleCount === 2 && femaleCount === 2 && uniqueRaces.size >= 3;
}

export function isInsideOut(group: Superhero[]): boolean {
  if (group.length !== 4) return false;
  const uniqueEmotions = new Set(group.map((hero) => hero.emotion));
  const uniqueRaces = new Set(group.map((hero) => hero.race));

  return uniqueEmotions.size === 4 && uniqueRaces.size >= 3;
}

export function isTheIncredibles(group: Superhero[]): boolean {
  if (group.length !== 4) return false;
  const uniqueAges = new Set(group.map((hero) => hero.age));
  const uniqueEmotions = new Set(group.map((hero) => hero.emotion));

  return uniqueAges.size === 4 && uniqueEmotions.size === 4;
}

// Main function to determine the group
export function determineHeroGroup(group: Superhero[]): number {
  if (isBalanceGuardians(group)) return 1;
  if (isInsideOut(group)) return 2;
  if (isTheIncredibles(group)) return 3;
  return -1;
}

// Function to generate all combinations of a certain size
export function combinations<T>(arr: T[], size: number): T[][] {
  const result: T[][] = [];
  const f = (start: number, combo: T[]) => {
    if (combo.length === size) {
      result.push(combo);
      return;
    }
    for (let i = start; i < arr.length; i++) {
      f(i + 1, combo.concat(arr[i]));
    }
  };
  f(0, []);
  return result;
}

// Group type definition
type GroupType = 'Balance Guardians' | 'Inside Out' | 'The Incredibles';

// Group result type definition
type GroupResult = {
  group: Superhero[];
  type: GroupType;
};

// Function to group heroes into valid groups
export function groupHeroes(heroes: Superhero[]): {
  groups: GroupResult[];
  points: number;
} {
  let points = 0;
  const groups: GroupResult[] = [];
  const usedFileNames = new Set<number>();

  const heroList = [...heroes];

  const tryFormGroup = (
    criteriaFn: (group: Superhero[]) => boolean,
    groupType: GroupType
  ): Superhero[] | null => {
    for (let i = 0; i < heroList.length - 3; i++) {
      if (usedFileNames.has(i)) continue;
      for (let j = i + 1; j < heroList.length - 2; j++) {
        if (usedFileNames.has(j)) continue;
        for (let k = j + 1; k < heroList.length - 1; k++) {
          if (usedFileNames.has(k)) continue;
          for (let l = k + 1; l < heroList.length; l++) {
            if (usedFileNames.has(l)) continue;
            const group = [heroList[i], heroList[j], heroList[k], heroList[l]];
            if (criteriaFn(group)) {
              usedFileNames.add(i);
              usedFileNames.add(j);
              usedFileNames.add(k);
              usedFileNames.add(l);
              groups.push({ group, type: groupType });
              points += 1;
              return group;
            }
          }
        }
      }
    }
    return null;
  };

  let group: Superhero[] | null;

  // Prioritize forming Balance Guardians
  while ((group = tryFormGroup(isBalanceGuardians, 'Balance Guardians'))) { console.log("still on balance") }

  // Then form Inside Out
  while ((group = tryFormGroup(isInsideOut, 'Inside Out'))) { console.log("still on inside out") }

  // Finally form The Incredibles
  while ((group = tryFormGroup(isTheIncredibles, 'The Incredibles'))) { console.log("still on incredibles") }

  return { groups, points };
}
