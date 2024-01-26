package main

const (
	COMMON Rarity = iota + 1
	FREE
	RARE
	EPIC
	LEGENDARY
)

const (
	GAME Type = iota + 1
	PLAYER
	HERO
	MINION
	SPELL
	ENCHANTMENT
	WEAPON
	ITEM
	TOKEN
	HERO_POWER
)

const (
	TEST_TEMPORARY Set = iota + 1
	BASIC
	EXPERT1
	HOF
	MISSIONS
	DEMO
	NONE
	CHEAT
	SetBLANK
	DEBUG_SP
	PROMO
	NAXX // Curse of Naxxramas
	GVG  // Goblins vs Gnomes
	BRM  // Blackrock Mountain
	TGT  // The Grand Tournament
	CREDITS
	HERO_SKINS
	TB // Tavern Brawl
	SLUSH
	LOE // The League of Explorers
	OG  // Whispers of the Old Gods
	OG_RESERVE
	KARA // One Night in Karazhan
	KARA_RESERVE
	GANGS // Mean Streets of Gadgetzan
	GANGS_RESERVE
	UNGORO                             // Journey to Un'Goro
	ICECROWN                Set = 1001 // Knights of the Frozen Throne
	LOOTAPALOOZA            Set = 1004 // Kobolds & Catacombs
	GILNEAS                 Set = 1125 // The Witchwood
	BOOMSDAY                Set = 1127 // The Boomsday Project
	TROLL                   Set = 1129 // Rastakhan's Rumble
	DALARAN                 Set = 1130 // Rise of Shadows
	ULDUM                   Set = 1158 // Saviours of Uldum
	DRAGONS                 Set = 1347 // Descent of Dragons
	YEAR_OF_THE_DRAGON      Set = 1403
	BLACK_TEMPLE            Set = 1414 // Ashes of Outlands
	WILD_EVENT              Set = 1439
	SCHOLOMANCE             Set = 1443 // Scholomance Academy
	BATTLEGROUNDS           Set = 1453
	DEMON_HUNTER_INITIATE   Set = 1463
	DARKMOON_FAIRE          Set = 1466 // Madness at the Darkmoon Faire
	THE_BARRENS             Set = 1525 // Forged in the Barrens
	WAILING_CAVERNS         Set = 1559
	STORMWIND               Set = 1578 // United in Stormwind
	LETTUCE                 Set = 1586 // Mercenaries
	ALTERAC_VALLEY          Set = 1626 // Fractured in Alterac Valley
	LEGACY                  Set = 1635
	CORE                    Set = 1637
	VANILLA                 Set = 1646
	THE_SUNKEN_CITY         Set = 1658 // Voyage to the Sunken City
	REVENDRETH              Set = 1691 // Murder at Castle Nathria
	MERCENARIES_DEV         Set = 1705
	RETURN_OF_THE_LICH_KING Set = 1776
	BATTLE_OF_THE_BANDS     Set = 1809
	TITANS                  Set = 1858
	PATH_OF_ARTHAS          Set = 1869
	WILD_WEST               Set = 1892
	WONDERS                 Set = 1898
	TUTORIAL                Set = 1904
)

const (
	DEATHKNIGHT Class = iota + 1
	DRUID
	HUNTER
	MAGE
	PALADIN
	PRIEST
	ROGUE
	SHAMAN
	WARLOCK
	WARRIOR
	DREAM
	NEUTRAL
	WHIZBANG
	DEMONHUNTER
)

func (t Type) String() string {
	switch t {
	case 0:
		return "invalid"
	case GAME:
		return "game"
	case PLAYER:
		return "player"
	case HERO:
		return "hero"
	case MINION:
		return "minion"
	case SPELL:
		return "spell"
	case ENCHANTMENT:
		return "enchantment"
	case WEAPON:
		return "weapon"
	case ITEM:
		return "item"
	case TOKEN:
		return "token"
	case HERO_POWER:
		return "hero power"
	default:
		return "unknown"
	}
}

func (r Rarity) String() string {
	switch r {
	case 0:
		return "invalid"
	case COMMON:
		return "common"
	case FREE:
		return "free"
	case RARE:
		return "rare"
	case EPIC:
		return "epic"
	case LEGENDARY:
		return "legendary"
	default:
		return "unknown"
	}
}

func (c Class) String() string {
	switch c {
	case 0:
		return "invalid"
	case DEATHKNIGHT:
		return "deathknight"
	case DRUID:
		return "druid"
	case HUNTER:
		return "hunter"
	case MAGE:
		return "mage"
	case PALADIN:
		return "paladin"
	case PRIEST:
		return "priest"
	case ROGUE:
		return "rogue"
	case SHAMAN:
		return "shaman"
	case WARLOCK:
		return "warlock"
	case WARRIOR:
		return "warrior"
	case DREAM:
		return "dream"
	case NEUTRAL:
		return "neutral"
	case WHIZBANG:
		return "whizbang"
	case DEMONHUNTER:
		return "demonhunter"
	default:
		return "unknown"
	}
}

func (c Set) String() string {
	switch c {
	case 0:
		return "invalid"
	case TEST_TEMPORARY:
		return "test_temporary"
	case BASIC:
		return "basic"
	case EXPERT1:
		return "expert1"
	case HOF:
		return "hof"
	case MISSIONS:
		return "missions"
	case DEMO:
		return "demo"
	case NONE:
		return "none"
	case CHEAT:
		return "cheat"
	case SetBLANK:
		return "blank"
	case DEBUG_SP:
		return "debug_sp"
	case PROMO:
		return "promo"
	case NAXX:
		return "Curse of Naxxramas"
	case GVG:
		return "Goblins vs Gnomes"
	case BRM:
		return "Blackrock Mountain"
	case TGT:
		return "The Grand Tournament"
	case CREDITS:
		return "credits"
	case HERO_SKINS:
		return "hero_skins"
	case TB:
		return "tavern brawl"
	case SLUSH:
		return "slush"
	case LOE:
		return "The League of Explorers"
	case OG:
		return "Whispers of the Old Gods"
	case OG_RESERVE:
		return "og_reserve"
	case KARA:
		return "One Night in Karazhan"
	case KARA_RESERVE:
		return "kara reserve"
	case GANGS:
		return "Mean Streets of Gadgetzan"
	case GANGS_RESERVE:
		return "gangs reserve"
	case UNGORO:
		return "Journey to Un'Goro"
	case ICECROWN:
		return "Knights of the Frozen Throne"
	case LOOTAPALOOZA:
		return "Kobolds & Catacombs"
	case GILNEAS:
		return "The Witchwood"
	case BOOMSDAY:
		return "The Boomsday Project"
	case TROLL:
		return "Rastakhan's Rumble"
	case DALARAN:
		return "Rise of Shadows"
	case ULDUM:
		return "Saviours of Uldum"
	case DRAGONS:
		return "Descent of Dragons"
	case YEAR_OF_THE_DRAGON:
		return "year of the dragon"
	case BLACK_TEMPLE:
		return "Ashes of Outlands"
	case WILD_EVENT:
		return "wild event"
	case SCHOLOMANCE:
		return "Scholomance Academy"
	case BATTLEGROUNDS:
		return "battlegrounds"
	case DEMON_HUNTER_INITIATE:
		return "demon hunter initiate"
	case DARKMOON_FAIRE:
		return "Madness at the Darkmoon Faire"
	case THE_BARRENS:
		return "Forged in the Barrens"
	case WAILING_CAVERNS:
		return "Wailing Caverns"
	case STORMWIND:
		return "United in Stormwind"
	case LETTUCE:
		return "Mercenaries"
	case ALTERAC_VALLEY:
		return "Fractured in Alterac Valley"
	case LEGACY:
		return "Legacy"
	case CORE:
		return "Core"
	case VANILLA:
		return "Vanilla"
	case THE_SUNKEN_CITY:
		return "Voyage to the Sunken City"
	case REVENDRETH:
		return "Murder at Castle Nathria"
	case MERCENARIES_DEV:
		return "MERCENARIES_DEV"
	case RETURN_OF_THE_LICH_KING:
		return "Return of the Lich King"
	case BATTLE_OF_THE_BANDS:
		return "Battle of the bands"
	case TITANS:
		return "Titans"
	case PATH_OF_ARTHAS:
		return "Path of Arthas"
	case WILD_WEST:
		return "wild west"
	case WONDERS:
		return "wonders"
	case TUTORIAL:
		return "tutorial"
	default:
		return "Unknown"
	}
}
