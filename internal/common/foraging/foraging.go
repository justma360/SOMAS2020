package foraging

// see https://colab.research.google.com/drive/1g1tiX27Ds7FGjj4_WjFB3OLj8Fat_Ur5?usp=sharing for experiments + simulations

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// defines the incremental increase in input resources required to move up a utility tier (to be able to hunt another deer)
var deerUtilityIncrements = []float64{1.0, 0.75, 0.5, 0.25} //TODO: move this to central config store

type DeerHuntParams struct {
	p float64 // Bernoulli p variable (whether or not a deer is caught)
	lam float64 // Exponential lambda (scale) param for W (weight variable)
}

// Island captures the location of a single island
type DeerHunt struct {
	participants  map[shared.ClientID][float64]
	params DeerHuntParams
}

func (d DeerHunt) totalInput float64 {
	i := 0.0
	for _, x := range d.participants {
		i+=x
	}
	return i
}

func (d DeerHunt) Hunt() float64 {
	input := d.totalInput()
	maxDeer := deerUtilityTier(input, deerUtilityIncrements) // get max number of deer allowed for given resource input
	utility := 0.0
	for i := 1; i < maxDeer; i++ {
		utility += deerReturn(d.params)
	}
	return utility
}

// deerUtilityTier gets the discrete utility tier (i.e. max number of deer) for given scalar input
func deerUtilityTier(input float64, increments []float64) int{
	if len(increments) == 0 || input < increments[0]{
		return 0
	}
	sum := increments[0]
	for i := 1; i < len(increments); i++ {
		if(input < sum){
			return i-1
		}
		sum+=increments[i]
    }
	return len(increments)
}

// deerReturn() is effectively the combination of two other RVs:
// - D: Bernoulli RV that represents the probaility of catching a deer at all (binary). Usually p - i.e. P(D=1) = p - will be fairly 
// close to 1 (fairly high chance of catching a deer if you invest the resources)
// - W: A continuous RV that adds some variance to the return. This could be interpreted as the weight of the deer that is caught. W is
// exponentially distributed such that the prevalence of deer of certain size is inversely prop. to the size.
// returns H, where H = D*(1+W) is an other random variable
func deerReturn(params DeerHuntParams) float64 {
	W := distuv.Exponential{Rate: params.lam} // Rate = lambda
	D := distuv.Bernoulli{P: params.p}        // Bernoulli RV where `P` = P(X=1)
	return D.Rand()*(1+W.Rand())
}