export type ElementKey = 'wood' | 'fire' | 'earth' | 'metal' | 'water';
export type TabKey = 'romance' | 'health' | 'career' | 'wealth';
export type Screen = 'onboard' | 'deep' | 'dashboard';

export const ELEMENTS: ElementKey[] = ['wood', 'fire', 'earth', 'metal', 'water'];

export type City = {
  name: string;
  country: string;
};

export const CITIES: City[] = [
  { name: 'Singapore', country: 'SGP' },
  { name: 'Hong Kong', country: 'HKG' },
  { name: 'Tokyo', country: 'JPN' },
  { name: 'Shanghai', country: 'CHN' },
  { name: 'Beijing', country: 'CHN' },
  { name: 'Taipei', country: 'TWN' },
  { name: 'Seoul', country: 'KOR' },
  { name: 'Bangkok', country: 'THA' },
  { name: 'Kuala Lumpur', country: 'MYS' },
  { name: 'Manila', country: 'PHL' },
  { name: 'Jakarta', country: 'IDN' },
  { name: 'Mumbai', country: 'IND' },
  { name: 'Delhi', country: 'IND' },
  { name: 'London', country: 'GBR' },
  { name: 'Paris', country: 'FRA' },
  { name: 'Berlin', country: 'DEU' },
  { name: 'Madrid', country: 'ESP' },
  { name: 'Rome', country: 'ITA' },
  { name: 'New York', country: 'USA' },
  { name: 'Los Angeles', country: 'USA' },
  { name: 'Chicago', country: 'USA' },
  { name: 'Toronto', country: 'CAN' },
  { name: 'Vancouver', country: 'CAN' },
  { name: 'Sydney', country: 'AUS' },
  { name: 'Melbourne', country: 'AUS' },
  { name: 'Auckland', country: 'NZL' },
  { name: 'Dubai', country: 'ARE' },
  { name: 'Cairo', country: 'EGY' },
  { name: 'Lagos', country: 'NGA' },
  { name: 'Nairobi', country: 'KEN' },
  { name: 'Johannesburg', country: 'ZAF' },
  { name: 'Sao Paulo', country: 'BRA' },
  { name: 'Mexico City', country: 'MEX' },
  { name: 'Buenos Aires', country: 'ARG' },
];

export type ElementMeta = {
  label: string;
  short: string;
  color: string;
};

export const ELEMENT_META: Record<ElementKey, ElementMeta> = {
  wood: { label: 'Wood', short: 'WOOD', color: '#5E8A4C' },
  fire: { label: 'Fire', short: 'FIRE', color: '#BD5B3B' },
  earth: { label: 'Earth', short: 'EARTH', color: '#A97C34' },
  metal: { label: 'Metal', short: 'METAL', color: '#6E7B8C' },
  water: { label: 'Water', short: 'WATER', color: '#35708E' },
};

export const ELEMENT_TEXT: Record<
  ElementKey,
  Record<TabKey, string>
> = {
  wood: {
    romance:
      "You love the way a tree loves — reaching outward, always growing the bond rather than settling into it. Partners feel your steady, expansive warmth, but you can outgrow connections that stop stretching with you.",
    health:
      "Vitality flows through your tendons and decision-making organs. Movement — real movement, not just intention — keeps your energy from turning brittle. Stagnation is your body's least favorite word.",
    career:
      "You thrive where there's room to expand: building, planting, initiating. Rigid hierarchies frustrate you; you do your best work when you're allowed to grow the role, not just fill it.",
    wealth:
      "Your money grows the way you do — slowly, then suddenly. You're better at cultivating long-term value than chasing quick wins, though impatience can tempt you to pull up the roots too early.",
  },
  fire: {
    romance:
      "You love loudly and visibly. Chemistry matters enormously to you, and you'd rather burn brightly for a short season than dim yourself for stability. Partners either match your heat or feel scorched by it.",
    health:
      "Your circulatory system and heart carry your signature. Passion fuels you, but unchecked intensity burns through your reserves fast — rest isn't optional, it's fuel.",
    career:
      "You're the spark in the room: charismatic, quick, magnetic to opportunity. You do best in visible, expressive roles, and worst in ones that ask you to stay quiet and wait your turn.",
    wealth:
      "Money moves fast around you — earned in bursts, spent with flair. Building lasting wealth means pairing your instinct for opportunity with someone (or something) that can bank the embers.",
  },
  earth: {
    romance:
      "You love like ground you can build a house on. Steady, dependable, unglamorous in the best way — you're the partner people come home to, not just the one they chase.",
    health:
      "Your digestive center is your barometer. When you're overextended caring for everyone else, it shows up first in your gut. Nourishment, including your own, is not indulgence.",
    career:
      "You're the stabilizer — the one who makes plans actually work. Reliable execution is your edge; just watch for being handed everyone else's unfinished business because you never say no.",
    wealth:
      "You save before you spend and plan before you leap. Your wealth grows through patience and diversification rather than bold bets — steady compounding is your natural mode.",
  },
  metal: {
    romance:
      "You love with precision — you notice details others miss and hold high standards for what a relationship should be. That clarity is a gift, but it can read as distance if you don't name what you feel.",
    health:
      "Your lungs and skin carry your tension first. Structure and clean air calm you; chaos and clutter — physical or emotional — drain you faster than almost anything else.",
    career:
      "You bring order to disorder. Editing, refining, systems, quality control — you make things sharper. You do poorly in loose, undefined roles with no clear standard to meet.",
    wealth:
      "You're disciplined with money to the point of austerity sometimes. You rarely overspend, and your instinct to cut what isn't working serves your portfolio as well as your closet.",
  },
  water: {
    romance:
      "You love the way water moves — adapting to whoever you're with, reading what they need before they say it. That flexibility draws people in, but you can lose your own shape if you're not careful.",
    health:
      "Your kidneys and adrenal reserves are your tell. You run on deep reserves rather than quick bursts, which means burnout, when it comes, comes from the bottom up. Protect your sleep.",
    career:
      "You're the strategist — comfortable with ambiguity, good at finding the path of least resistance to a goal. You underperform in rigid, script-driven roles that don't let you improvise.",
    wealth:
      "Your money finds unconventional channels — you're drawn to opportunities others overlook. The risk is spreading too thin; depth in a few currents beats width across all of them.",
  },
};

export const LIFE_PATH_TEXT: Record<number, string> = {
  1: "an independent starter — you'd rather build the first version yourself than wait for consensus.",
  2: "a natural mediator — you sense what a room needs before it's said out loud.",
  3: "expressive by design — your ideas want an audience, not just a notebook.",
  4: "built for the long haul — structure isn't a constraint to you, it's how things get real.",
  5: "wired for change — routine is the fastest way to lose your interest entirely.",
  6: "responsible to your core — people orbit you because you actually show up.",
  7: "an investigator at heart — you trust what you've examined over what you've been told.",
  8: "oriented toward scale — you think in systems, leverage, and long games.",
  9: "finishing what others start — you see the whole arc before most people see the first act.",
  11: "intuitive to an unusual degree — you often know before you can explain why.",
  22: "a builder of lasting things — your ideas are rarely small in scope.",
  33: "guided by care for others — your gift is making people feel genuinely seen.",
};
