//
//  WHONewMessageViewController.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "WHONewMessageViewController.h"

@interface WHONewMessageViewController ()

@end

@implementation WHONewMessageViewController

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self) {
        // Custom initialization
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    UIColor* hooliColor = [UIColor colorWithRed:109.0/255 green:211.0/255 blue:170.0/255 alpha:1.0];
    UIColor* brownColor = [UIColor colorWithRed:78.0/255 green:46.0/255 blue:40.0/255 alpha:1.0];
    UIColor* goldColor = [UIColor colorWithRed:198.0/255 green:150.0/255 blue:73.0/255 alpha:1.0];
//    UIColor* brickColor = [UIColor colorWithRed:207.0/255 green:86.0/255 blue:61.0/255 alpha:1.0];
    
    self.messageField.delegate = self;
    self.messageField.textAlignment = NSTextAlignmentCenter;
    self.messageField.font = [UIFont fontWithName:@"AvenirNextCondensed-DemiBold" size:17.0];
    
//    [self.navigationItem.backBarButtonItem setTintColor:brownColor];
    [self.navigationController.navigationBar setTintColor:brownColor];
    
    UILabel* titleLabel = [[UILabel alloc] init];
    [titleLabel setText:@"Hooli"];
    [titleLabel setFont:[UIFont fontWithName:@"Superclarendon-BlackItalic" size:25.0]];
    [titleLabel setTextColor:goldColor];
//        [titleLabel setAlpha:0.75];
    [titleLabel.layer setShadowColor:[UIColor darkGrayColor].CGColor];
    [titleLabel.layer setShadowOffset:(CGSize) { .width = 1.0, .height = 1.0 }];
    [titleLabel.layer setShadowRadius:1.0];
    [titleLabel.layer setShadowOpacity:.65];
    [titleLabel sizeToFit];
    [self.navigationItem setTitleView:titleLabel];
    
    [self.view setBackgroundColor:brownColor];
    [self.messageField setBackgroundColor:brownColor];
    [self.messageField setTextColor:hooliColor];
    
    [self.messageField addObserver:self forKeyPath:@"contentSize" options:(NSKeyValueObservingOptionNew) context:NULL];
    
    //for Done button on keyboard
//    [self.messageField addTarget:self action:@selector(messageView:) forControlEvents:UIControlEventEditingDidEndOnExit];
}

//center the text vertically in the UITextView
-(void)observeValueForKeyPath:(NSString *)keyPath ofObject:(id)object change:(NSDictionary *)change context:(void *)context {
    UITextView *tv = object;
    CGFloat topCorrect = ([tv bounds].size.height - [tv contentSize].height * [tv zoomScale])/2.0;
    topCorrect = ( topCorrect < 0.0 ? 0.0 : topCorrect );
    tv.contentOffset = (CGPoint){.x = 0, .y = -topCorrect};
}

- (void)submit:(id)sender {
    [self.delegate receivedNewMessage:self.messageField.text];
    [UIView animateWithDuration:0.7
                     animations:^{
                         [UIView setAnimationCurve:UIViewAnimationCurveEaseInOut];
                         [UIView setAnimationTransition:UIViewAnimationTransitionCurlDown forView:self.navigationController.view cache:NO];
                     }];
    [self.navigationController popViewControllerAnimated:YES];
}

- (void)showPostButton {
    UIColor* brownColor = [UIColor colorWithRed:78.0/255 green:46.0/255 blue:40.0/255 alpha:1.0];
    UIBarButtonItem* submitButton = [[UIBarButtonItem alloc] initWithTitle:@"Post" style:UIBarButtonItemStylePlain target:self action:@selector(submit:)];
    [submitButton setTintColor:brownColor];
    [self.navigationItem setRightBarButtonItem:submitButton animated:YES];
}

- (void)hidePostButton {
    [self.navigationItem setRightBarButtonItem:nil animated:YES];
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (IBAction)messageView:(id)sender {
    [self.messageField resignFirstResponder];
}

-(void)textViewDidBeginEditing:(UITextView *)textView {
    [UIView transitionWithView:self.placeholderLabel
                      duration:.3
                       options:UIViewAnimationOptionTransitionCrossDissolve
                    animations:NULL
                    completion:NULL];
    self.placeholderLabel.hidden = YES;
}

-(void)textViewDidChange:(UITextView *)textView {
//    self.placeholderLabel.hidden = [self.messageField.text length] > 0;
    if ([self.messageField.text length] > 0) {
        if (self.navigationItem.rightBarButtonItem == nil) {
            [self showPostButton];
        }
    }
    else if (self.navigationItem.rightBarButtonItem != nil) {
        [self hidePostButton];
    }
}

-(void)textViewDidEndEditing:(UITextView *)textView {
    [UIView transitionWithView:self.placeholderLabel
                      duration:.3
                       options:UIViewAnimationOptionTransitionCrossDissolve
                    animations:NULL
                    completion:NULL];
    self.placeholderLabel.hidden = [self.messageField.text length] > 0;
}

//for Done button
- (BOOL)textView:(UITextView *)textView shouldChangeTextInRange:(NSRange)range replacementText:(NSString *)text {
    if([text isEqualToString:@"\n"]) {
        [textView resignFirstResponder];
        return NO;
    }
    return YES;
}

@end
