import fs from 'fs';
import crypto from 'crypto';
import {
  Age,
  Emotion,
  Gender,
  getRandomEnumValue,
  groupHeroes,
  Race,
  readSuperheroesFromFile,
  type Superhero,
} from './helper';

// Generate 100 random superheroes
const superheroes: Superhero[] = Array.from({ length: 1500 }, () => ({
  uuid: crypto.randomUUID(),
  gender: getRandomEnumValue(Gender),
  emotion: getRandomEnumValue(Emotion),
  age: getRandomEnumValue(Age),
  race: getRandomEnumValue(Race),
}));

// Write to JSON file
// fs.writeFileSync(
//   'superheroes.json',
//   JSON.stringify(superheroes, null, 2),
//   'utf-8'
// );

// Function to read and parse the JSON file

const heroes = readSuperheroesFromFile('superheroes.json');
// console.log(heroes);
fs.writeFile(
  'combination.json',
  JSON.stringify(groupHeroes(superheroes), null, 2),
  (err) => {
    console.log(err);
  }
);
