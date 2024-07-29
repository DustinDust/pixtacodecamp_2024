import input from './superheroes.json';
import fs from 'fs';


const filtered = input.filter((hero) => {
    if (hero.age === "20-30s" && hero.emotion === "Happiness" || hero.age === "20-30s" && hero.emotion == "Neutral") {
        return true
    }
})

const minor_filter = input.filter(hero => {
    if (hero.age === "20-30s" && hero.emotion === "Happiness" || hero.age === "20-30s" && hero.emotion == "Neutral") {
        return false
    }
    return true
})
console.log(filtered.length);
console.log(minor_filter.length);

fs.writeFileSync("major_input.json", JSON.stringify(filtered, null, 2))
fs.writeFileSync("minor_input.json", JSON.stringify(minor_filter, null, 2))
