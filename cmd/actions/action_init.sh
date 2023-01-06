if [[ "$PROMPT_COMMAND" == *__tdc_ps1* ]]; then
	exit 0
fi

TCC_TIME=0

__tdc_ps1() {	
    local tcc_elapsed=$((SECONDS - TCC_TIME))
    if [[ "$TCC_TIME" -eq "0" || "$tcc_elapsed" -gt "#PERIOD#" ]]; then
        TCC_TIME=$SECONDS		
        #THISPATH# notify
    fi
}

PROMPT_COMMAND="__tdc_ps1 $PROMPT_COMMAND"
alias todo=#THISPATH# 
echo "todo-cli ps1 initialized (every #PERIOD# seconds)"