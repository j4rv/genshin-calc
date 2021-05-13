vueApp.component('char-build', {
    data: () => ({
        HP: 10000,
        ATK: 1000,
        DEF: 800,
        CRIT: 5,
        CDMG: 50,
        DMGP: 0,
        ER: 100,
        atkMV: 100,
        hpMV: 0,
        defMV: 0,
    }),
    methods: {
        copyCharBuild() {
            let state = this.$data
            let jsonState = JSON.stringify(state)
            localStorage.setItem("char-build", jsonState)
        },
        pasteCharBuild() {
            let jsonState = localStorage.getItem("char-build")
            let state = JSON.parse(jsonState)
            for (key in state) {
                if (["savedBuilds"].includes(key)) {
                    continue
                }
                this[key] = state[key]
            }
        },
    },
    computed: {
        CritAvgM: function () { return 1 + Math.min(this.CRIT, 100) * this.CDMG / 10000 },
        finalDMG: function () {
            return (
                (1 + this.DMGP / 100) * (
                    this.atkMV / 100 * this.ATK +
                    this.hpMV / 100 * this.HP +
                    this.defMV / 100 * this.DEF
                )
            )
        },
    },
    template: `
    <div class="inputSet" style="margin: 16px;">
        <div>
            <input id="Notes" style="width: 12em;" type="string" value="Notes">
        </div>
        <h3>Stats <a target="_blank"
                href="https://genshin-impact.fandom.com/wiki/Attributes#Base_Stats">❓</a>
        </h3>
        <div>
            <label for="HP">HP</label>
            <input id="HP" v-model.number="HP" type="number">
        </div>
        <div>
            <label for="ATK">ATK</label>
            <input id="ATK" v-model.number="ATK" type="number">
        </div>
        <div>
            <label for="DEF">DEF</label>
            <input id="DEF" v-model.number="DEF" type="number">
        </div>
        <div>
            <label for="CRIT">Base Crit%</label>
            <input id="CRIT" v-model.number="CRIT" type="number" step=".1">
        </div>
        <div>
            <label for="CDMG">Base Crit DMG%</label>
            <input id="CDMG" v-model.number="CDMG" type="number" step=".1">
        </div>
        <div>
            <label for="ER">Energy recharge</label>
            <input id="ER" v-model.number="ER" type="number" step=".1">
        </div>
        <div>
            <label for="DMGP">Base Damage %</label>
            <input id="DMGP" v-model.number="DMGP" type="number" step=".1">
        </div>

        <h3>MVs</h3>
        <div>
            <label for="atkMV">ATK Motion value (%)</label>
            <input id="atkMV" v-model.number="atkMV" type="number" step=".1">
        </div>
        <div>
            <label for="hpMV">HP Motion value (%) (Ex: Zhongli's Q)</label>
            <input id="hpMV" v-model.number="hpMV" type="number" step=".01">
        </div>
        <div>
            <label for="defMV">DEF Motion value (%) (Ex: Albedo's E)</label>
            <input id="defMV" v-model.number="defMV" type="number" step=".01">
        </div>

        <h3>Damage output</h3>
        <table>
            <tr>
                <th>Non-crit DMG</th>
                <td>{{ Math.trunc(finalDMG) }}</td>
            </tr>
            <tr>
                <th>Crit DMG</th>
                <td>{{ Math.trunc(finalDMG * ( 1 + CDMG / 100)) }}</td>
            </tr>
            <tr>
                <th>Average DMG</th>
                <td>{{ Math.trunc(finalDMG * CritAvgM) }}</td>
            </tr>
        </table>

        <button @click="copyCharBuild()">Copy</button>
        <button @click="pasteCharBuild()">Paste</button>
    </div>
  `
})
