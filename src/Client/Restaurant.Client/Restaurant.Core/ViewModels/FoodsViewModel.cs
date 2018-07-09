﻿using System;
using System.Collections.ObjectModel;
using System.Reactive.Linq;
using System.Threading.Tasks;
using System.Windows.Input;
using ReactiveUI;
using Restaurant.Abstractions.Adapters;
using Restaurant.Abstractions.Api;
using Restaurant.Abstractions.Factories;
using Restaurant.Abstractions.Services;
using Restaurant.Abstractions.ViewModels;
using Restaurant.Common.Constants;
using Restaurant.Common.DataTransferObjects;

namespace Restaurant.Core.ViewModels
{
    public class FoodsViewModel : BaseViewModel
    {
        private readonly IFoodDetailViewModelFactory _foodDetailViewModelFactory;
        private readonly IFoodsApi _foodsApi;
        private readonly INavigationService _navigationService;
        private ObservableCollection<FoodDto> _foods;
        private FoodDto _selectedFood;

	    internal FoodsViewModel()
	    {   
	    }

        public FoodsViewModel(
            IBasketViewModel basketViewModel,
            IFoodsApi foodsApi,
            INavigationService navigationService,
            IFoodDetailViewModelFactory foodDetailViewModelFactory)
        {
            _foodsApi = foodsApi;
            _navigationService = navigationService;
            _foodDetailViewModelFactory = foodDetailViewModelFactory;

            this.WhenAnyValue(x => x.SelectedFood)
                .Where(x => x != null)
                .Subscribe(async food => await NavigateToFoodDetail(food));

            GoToBasket =
                ReactiveCommand.CreateFromTask(async () => await _navigationService.NavigateAsync(BasketViewModel));

            BasketViewModel = basketViewModel;
        }

        public ObservableCollection<FoodDto> Foods
        {
            get => _foods;
            private set => this.RaiseAndSetIfChanged(ref _foods, value);
        }

        public FoodDto SelectedFood
        {
            get => _selectedFood;
            set => this.RaiseAndSetIfChanged(ref _selectedFood, value);
        }

        public IBasketViewModel BasketViewModel { get; }

        public ICommand GoToBasket { get; set; }

        public override string Title => "Foods";

        public async Task LoadFoods()
        {
            IsLoading = true;
            var foods = await _foodsApi.GetFoods(5, 0);
            if (!BootstrapperBase.MockData)
            {
                foreach (var food in foods)
                    food.Picture = ApiConstants.ApiClientUrl + food.Picture;
            }
            Foods = new ObservableCollection<FoodDto>(foods);
            BasketViewModel.RaiseOrdersCount();
            IsLoading = false;
        }

        private async Task NavigateToFoodDetail(FoodDto food)
        {
            var viewModel = _foodDetailViewModelFactory.GetFoodDetailViewModel(food);
            await _navigationService.NavigateAsync(viewModel);
        }
    }
}