import type { ClassEntry, SpeciesEntry, BackgroundEntry, SpellEntry, FeatEntry, FeatureEntry } from './types';

export const mockClasses: ClassEntry[] = [
  {
    id: "fighter",
    name: "Fighter",
    hitDie: 10,
    spellcaster: false,
    primaryAbility: ["STR", "DEX"],
    savingThrows: ["STR", "CON"],
    subclassLevel: 3,
    subClasses: [
      { id: "battle-master", name: "Battle Master", description: "Combat maneuvers and superiority dice." },
      { id: "champion", name: "Champion", description: "Improved critical hits and regeneration." },
      { id: "eldritch-knight", name: "Eldritch Knight", description: "Evocation and abjuration magic." }
    ],
    features: [
      { name: "Second Wind", level: 1, description: "Regain 1d10 + level HP as a bonus action." },
      { name: "Fighting Style", level: 1, description: "Choose a fighting style." },
      { name: "Action Surge", level: 2, description: "Take an additional action on your turn." },
      { name: "Weapon Mastery", level: 1, description: "Mastery with two simple weapons." }
    ],
    spellcasting: null,
    skillChoices: { count: 2, from: ["acrobatics", "animal-handling", "athletics", "history", "insight", "intimidation", "perception", "survival"] }
  },
  {
    id: "wizard",
    name: "Wizard",
    hitDie: 6,
    spellcaster: true,
    primaryAbility: ["INT"],
    savingThrows: ["INT", "WIS"],
    subclassLevel: 3,
    subClasses: [
      { id: "abjuration", name: "Abjuration", description: "Protection and magical barriers." },
      { id: "evocation", name: "Evocation", description: "Destructive damage magic." },
      { id: "divination", name: "Divination", description: "Glimpses of the future." }
    ],
    features: [
      { name: "Arcane Recovery", level: 1, description: "Recover spell slots on a short rest." },
      { name: "Ritual Casting", level: 1, description: "Cast ritual spells from spellbook." },
      { name: "Spell Mastery", level: 5, description: "Cast 1st and 2nd level spells at will." }
    ],
    spellcasting: {
      ability: "INT",
      preparedSpells: [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
      knownSpells: false,
      slots: {
        "1": [2,3,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4],
        "2": [0,0,2,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3],
        "3": [0,0,0,0,2,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3]
      }
    },
    skillChoices: { count: 2, from: ["arcana", "history", "insight", "investigation", "medicine", "religion"] }
  }
];

export const mockSpecies: SpeciesEntry[] = [
  {
    id: "human",
    name: "Human",
    description: "Versatile and ambitious. Adapt to any role.",
    traits: [
      { name: "Versatile", description: "+1 to all ability scores." },
      { name: "Extra Language", description: "Learn one extra language." }
    ],
    abilityBonuses: { STR: 1, DEX: 1, CON: 1, INT: 1, WIS: 1, CHA: 1 },
    size: "Medium",
    speed: 30,
    languages: ["Common"],
    variants: []
  },
  {
    id: "elf",
    name: "Elf",
    description: "Graceful and long-lived. Deep connection to nature and magic.",
    traits: [
      { name: "Darkvision", description: "See in darkness up to 60 feet." },
      { name: "Trance", description: "Rest in 4 hours instead of 8." },
      { name: "Fey Ancestry", description: "Immunity to charm, advantage against being charmed." }
    ],
    abilityBonuses: { DEX: 2 },
    size: "Medium",
    speed: 30,
    languages: ["Common", "Elvish"],
    variants: [
      { id: "high-elf", name: "High Elf", description: "+1 INT, extra cantrip." },
      { id: "wood-elf", name: "Wood Elf", description: "+1 WIS, speed 35, mask of the wild." }
    ]
  },
  {
    id: "dwarf",
    name: "Dwarf",
    description: "Resilient and determined. Masters of stone and metal.",
    traits: [
      { name: "Darkvision", description: "See in darkness up to 60 feet." },
      { name: "Dwarven Resilience", description: "Resistance to poison." },
      { name: "Stonecunning", description: "Proficiency in stone history." }
    ],
    abilityBonuses: { CON: 2 },
    size: "Medium",
    speed: 25,
    languages: ["Common", "Dwarvish"],
    variants: [
      { id: "hill-dwarf", name: "Hill Dwarf", description: "+1 WIS, max HP +1/level." },
      { id: "mountain-dwarf", name: "Mountain Dwarf", description: "+2 STR, proficiency with light and medium armor." }
    ]
  },
  {
    id: "halfling",
    name: "Halfling",
    description: "Small but courageous. Natural luck and agility.",
    traits: [
      { name: "Lucky", description: "Reroll 1s on attacks, checks, or saves." },
      { name: "Brave", description: "Advantage against fear." },
      { name: "Halfling Nimbleness", description: "Move through larger creatures' spaces." }
    ],
    abilityBonuses: { DEX: 2 },
    size: "Small",
    speed: 25,
    languages: ["Common", "Halfling"],
    variants: [
      { id: "lightfoot", name: "Lightfoot", description: "+1 CHA, hide behind larger creatures." },
      { id: "stout", name: "Stout", description: "+1 CON, resistance to poison." }
    ]
  }
];

export const mockBackgrounds: BackgroundEntry[] = [
  {
    id: "acolyte",
    name: "Acolyte",
    description: "You served at a temple or shrine. You know rituals and dogmas.",
    skillProficiencies: ["insight", "religion"],
    toolProficiencies: [],
    languages: [2],
    equipment: ["Holy symbol", "Prayer book", "5 sticks of incense", "Vestments", "Common clothes", "15 gp"],
    feature: { name: "Shelter of the Faithful", description: "Temples of your god provide shelter and care." }
  },
  {
    id: "sage",
    name: "Sage",
    description: "You dedicated your life to study. You accumulated knowledge in libraries and academies.",
    skillProficiencies: ["arcana", "history"],
    toolProficiencies: [],
    languages: [2],
    equipment: ["Bottle of black ink", "Quill", "Small knife", "Letter from a dead colleague", "Common clothes", "10 gp"],
    feature: { name: "Researcher", description: "You know where to find information — libraries, shrines, sages." }
  },
  {
    id: "soldier",
    name: "Soldier",
    description: "You served in an army or militia. You know hierarchy and combat.",
    skillProficiencies: ["athletics", "intimidation"],
    toolProficiencies: ["gaming-set"],
    languages: [0],
    equipment: ["Insignia of rank", "Trophy from a fallen enemy", "Bone dice or deck of cards", "Common clothes", "10 gp"],
    feature: { name: "Military Rank", description: "Loyal soldiers recognize your authority and provide assistance." }
  },
  {
    id: "criminal",
    name: "Criminal",
    description: "You lived on the wrong side of the law. You know the underworld and its secrets.",
    skillProficiencies: ["deception", "stealth"],
    toolProficiencies: ["thieves-tools", "gaming-set"],
    languages: [0],
    equipment: ["Crowbar", "Dark common clothes with hood", "belt pouch", "15 gp"],
    feature: { name: "Criminal Contact", description: "You have a reliable contact in the criminal underworld." }
  }
];

export const mockSpells: SpellEntry[] = [
  { id: "fire-bolt", name: "Fire Bolt", level: 0, school: "Evocation", description: "Ranged fire attack. 1d10 damage." },
  { id: "mage-hand", name: "Mage Hand", level: 0, school: "Conjuration", description: "Spectral hand that manipulates objects." },
  { id: "prestidigitation", name: "Prestidigitation", level: 0, school: "Transmutation", description: "Minor magical tricks." },
  { id: "shield", name: "Shield", level: 1, school: "Abjuration", description: "+5 AC until the start of your next turn." },
  { id: "magic-missile", name: "Magic Missile", level: 1, school: "Evocation", description: "3 darts of force that automatically hit." },
  { id: "burning-hands", name: "Burning Hands", level: 1, school: "Evocation", description: "Cone of fire. 3d6 damage." },
  { id: "charm-person", name: "Charm Person", level: 1, school: "Enchantment", description: "The creature becomes friendly." },
  { id: "detect-magic", name: "Detect Magic", level: 1, school: "Divination", description: "Sense magic within 30 feet." },
  { id: "misty-step", name: "Misty Step", level: 2, school: "Conjuration", description: "Teleport up to 30 feet." },
  { id: "scorching-ray", name: "Scorching Ray", level: 2, school: "Evocation", description: "3 rays of fire. 2d6 each." }
];

export const mockFeats: FeatEntry[] = [
  { id: "alert", name: "Alert", prerequisites: {}, description: "+5 initiative, can't be surprised, no advantage on attacks from hidden enemies." },
  { id: "war-caster", name: "War Caster", prerequisites: { spellcasting: true }, description: "Advantage on concentration saves, cast spells as opportunity attacks." },
  { id: "tough", name: "Tough", prerequisites: {}, description: "HP max increases by 2 per level." },
  { id: "lucky", name: "Lucky", prerequisites: {}, description: "3 luck points to reroll attacks, checks, or saves." }
];

export const mockFeatures: FeatureEntry[] = [
  { id: "second-wind", classId: "fighter", name: "Second Wind", level: 1, description: "Regain 1d10 + level HP as a bonus action." },
  { name: "Fighting Style", level: 1, description: "Choose a fighting style." },
  { name: "Action Surge", level: 2, description: "Take an additional action on your turn." },
  { name: "Weapon Mastery", level: 1, description: "Mastery with two simple weapons." },
  { id: "arcane-recovery", classId: "wizard", name: "Arcane Recovery", level: 1, description: "Recover spell slots on a short rest." },
  { name: "Ritual Casting", level: 1, description: "Cast ritual spells from spellbook." },
  { name: "Spell Mastery", level: 5, description: "Cast 1st and 2nd level spells at will." }
];